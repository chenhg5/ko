package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	kitjwt "github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/etcdv3"
	"github.com/go-kit/kit/sd/lb"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Api设计关键类，restful/json转换/数据格式指定/错误编码列表
// 数据转换
// POST/PUT/DELETE/GET

// 请求与响应关系
// client   ->    gatewayProxy   ->    svc   ->   gatewayProxy  ->   client
//                   server          server          client
//                  decodeReq                       encodeReq
//                  encodeRes                       decodeRes
//
// 请求流
// gatewayProxy.server.decodeReq -> gatewayProxy.client.encodeReq -> gatewayProxy.client.decodeRes -> gatewayProxy.server.encodeRes

// 服务工厂生成器
func SvcFactory(method string, path string) sd.Factory {
	return func(instance string) (endpoint.Endpoint, io.Closer, error) {
		if !strings.HasPrefix(instance, "http") {
			instance = "http://" + instance
		}
		tgt, err := url.Parse(instance)
		fmt.Println("listening svc: ", tgt)
		if err != nil {
			return nil, nil, err
		}
		tgt.Path = path

		var (
			enc httptransport.EncodeRequestFunc
			dec httptransport.DecodeResponseFunc
		)

		method = strings.ToUpper(method)

		if method == "GET" {
			enc, dec = EncodeGetRequest, DecodeGetResponse
		}
		if method == "POST" {
			enc, dec = EncodeJsonRequest, DecodeGetResponse
		}

		return httptransport.NewClient(method, tgt, enc, dec).Endpoint(), nil, nil
	}
}

func MakeHandler(
	logger log.Logger,
	ins *etcdv3.Instancer,
	method string,
	path string,
	middlewares ...endpoint.Middleware,
) *httptransport.Server {
	factory := SvcFactory(method, path)
	endpointer := sd.NewEndpointer(ins, factory, logger)
	balancer := lb.NewRoundRobin(endpointer)
	retry := lb.Retry(3, 3*time.Second, balancer)

	for _, middleware := range middlewares {
		retry = middleware(retry)
	}

	opts := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
		httptransport.ServerErrorEncoder(encodeError),
	}

	if method == "POST" {
		return httptransport.NewServer(retry, DecodeJsonRequest, EncodeJSONResponse, opts...)
	} else {
		return httptransport.NewServer(retry, DecodeGetRequest, EncodeJSONResponse, opts...)
	}
}

func MakeJwtHandler(
	logger log.Logger,
	ins *etcdv3.Instancer,
	method string,
	path string,
	middlewares ...endpoint.Middleware,
) *httptransport.Server {
	factory := SvcFactory(method, path)
	endpointer := sd.NewEndpointer(ins, factory, logger)
	balancer := lb.NewRoundRobin(endpointer)
	retry := lb.Retry(3, 3*time.Second, balancer)

	for _, middleware := range middlewares {
		retry = middleware(retry)
	}

	opts := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
		httptransport.ServerErrorEncoder(encodeError),
		httptransport.ServerBefore(kitjwt.HTTPToContext()), // jwt token
	}

	if method == "POST" {
		return httptransport.NewServer(retry, DecodeJsonRequest, EncodeJSONResponse, opts...)
	} else {
		return httptransport.NewServer(retry, DecodeGetRequest, EncodeJSONResponse, opts...)
	}
}

// encode errors from business-logic
// svc panic
func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": http.StatusInternalServerError,
		"msg":  err.Error(),
	})
}

type Router struct {
	r      *mux.Router
	ins    *etcdv3.Instancer
	logger log.Logger
}

func InitRouter(logger log.Logger) *Router {
	var router Router
	router.r = mux.NewRouter()
	router.logger = logger
	return &router
}

func (router *Router) Service(prefix string, etcdClient etcdv3.Client) {
	router.ins = GetEtcdInstancer(prefix, etcdClient, router.logger)
}

func (router *Router) Post(path string, middlewares ...endpoint.Middleware) {
	router.r.Handle(path, MakeHandler(
		router.logger,
		router.ins,
		"POST",
		path,
		middlewares...,
	))
}

func (router *Router) Get(path string, middlewares ...endpoint.Middleware) {
	router.r.Handle(path, MakeHandler(
		router.logger,
		router.ins,
		"GET",
		path,
		middlewares...,
	))
}

func (router *Router) JwtPost(path string, middlewares ...endpoint.Middleware) {
	router.r.Handle(path, MakeJwtHandler(
		router.logger,
		router.ins,
		"POST",
		path,
		middlewares...,
	))
}

func (router *Router) JwtGet(path string, middlewares ...endpoint.Middleware) {
	router.r.Handle(path, MakeJwtHandler(
		router.logger,
		router.ins,
		"GET",
		path,
		middlewares...,
	))
}

// 客户端到内部服务：转换Get请求
func EncodeGetRequest(_ context.Context, req *http.Request, request interface{}) error {
	data := request.(commonUrlReq)
	req.URL.Path = strings.Replace(req.URL.Path, "{param}", data.Param, -1)

	return nil
}

// 客户端到内部服务：转换Json请求
func EncodeJsonRequest(_ context.Context, req *http.Request, request interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	req.Body = ioutil.NopCloser(&buf)

	return nil
}

// 客户端到内部服务：转换Json响应
func EncodeJSONResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

// 内部服务到客户端：解码Get响应
func DecodeGetResponse(ctx context.Context, resp *http.Response) (interface{}, error) {
	var commonResponse commonRes
	var outputResponse outputRes
	if err := json.NewDecoder(resp.Body).Decode(&commonResponse); err != nil {
		return nil, err
	}
	if commonResponse.Err != "" {
		outputResponse.Msg = commonResponse.Err
		outputResponse.Code = 500
		outputResponse.Data = map[string]interface{}{}
	} else {
		outputResponse.Msg = commonResponse.Msg
		outputResponse.Code = commonResponse.Code
		outputResponse.Data = commonResponse.Data
	}
	return outputResponse, nil
}

// 内部服务到客户端：解码Get请求
func DecodeGetRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	vars := mux.Vars(req)
	param, err := vars["param"]

	if !err {
		return nil, errBadRoute
	}
	var getReq commonUrlReq
	getReq.Param = param
	return getReq, nil
}

// 内部服务到客户端：解码Json请求
func DecodeJsonRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	var request commonJsonReq
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

// 全局请求与响应类型定义
type commonJsonReq struct {
	Param map[string]interface{} `json:"param"`
}
type commonUrlReq struct {
	Param string `json:"param"`
}
type commonRes struct {
	Code int                    `json:"code"`
	Msg  string                 `json:"msg"`
	Data map[string]interface{} `json:"data"`
	Err  string                 `json:"err,omitempty"`
}
type outputRes struct {
	Code int                    `json:"code"`
	Msg  string                 `json:"msg"`
	Data map[string]interface{} `json:"data"`
}

// 错误码
var errBadRoute = errors.New("10111 错误的路由参数")
