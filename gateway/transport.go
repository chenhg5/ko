package gateway

import (
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/endpoint"
	"io"
	"strings"
	"net/url"
	"github.com/gorilla/mux"
	"context"
	httptransport "github.com/go-kit/kit/transport/http"
	"net/http"
	"encoding/json"
	"errors"
	"bytes"
	"io/ioutil"
	"fmt"
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
func SvcFactory(ctx context.Context, method, path string) sd.Factory {
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
		if method == "GET" {
			enc, dec = EncodeGetRequest, DecodeGetResponse
		}
		if method == "POST" {
			enc, dec = EncodeJsonRequest, DecodeGetResponse
		}

		return httptransport.NewClient(method, tgt, enc, dec).Endpoint(), nil, nil
	}
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
	var getuserResponse commonRes
	getuserResponse.Msg  = "ok"
	getuserResponse.Code = 0
	if err := json.NewDecoder(resp.Body).Decode(&getuserResponse); err != nil {
		return nil, err
	}
	return getuserResponse, nil
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
	Param  map[string]interface{}   `json:"param"`
}
type commonUrlReq struct {
	Param string `json:"param"`
}
type commonRes struct {
	Code  int                      `json:"code"`
	Msg   string                   `json:"msg"`
	Data  map[string]interface{}   `json:"data"`
	Err   error                    `json:"err,omitempty"`
}

// 错误码
var errBadRoute = errors.New("10111 错误的路由参数")