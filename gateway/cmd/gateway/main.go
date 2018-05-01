package main

import (
	"time"
	"github.com/go-kit/kit/sd/etcdv3"
	"context"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/lb"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/endpoint"
	"io"
	"github.com/gorilla/mux"
	httptransport "github.com/go-kit/kit/transport/http"
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"fmt"
	"strings"
	"net/url"
	"errors"
)

func main() {
	// 1. 配置
	var (
		etcdServer = "localhost:2379"      // in the change from v2 to v3, the schema is no longer necessary if connecting directly to an etcd v3 instance
		ctx        = context.Background()
		httpAddr   = ":4001"
	)

	options := etcdv3.ClientOptions{
		// Path to trusted ca file
		CACert: "",

		// Path to certificate
		Cert: "",

		// Path to private key
		Key: "",

		// Username if required
		Username: "",

		// Password if required
		Password: "",

		// If DialTimeout is 0, it defaults to 3s
		DialTimeout: time.Second * 3,

		// If DialKeepAlive is 0, it defaults to 3s
		DialKeepAlive: time.Second * 3,
	}

	// 2. 服务发现
	etcdClient, err := etcdv3.NewClient(ctx, []string{etcdServer}, options)
	if err != nil {
		panic(err)
	}

	// 1) 用户中心服务

	// 创造实例
	barPrefix := "/svc/ucenter"
	logger := log.NewNopLogger()
	instancer, err := etcdv3.NewInstancer(etcdClient, barPrefix, logger)
	if err != nil {
		panic(err)
	}

	// 路由控制器构造
	factory := svcFactory(ctx, "GET", "/svc/ucenter/v1/user/{UID}")
	endpointer := sd.NewEndpointer(instancer, factory, logger)
	balancer := lb.NewRoundRobin(endpointer)
	retry := lb.Retry(3, 3*time.Second, balancer)

	// 路由
	r := mux.NewRouter()
	r.Handle("/svc/ucenter/v1/user/{UID}", httptransport.NewServer(retry, decodeGetRequest, encodeJSONResponse))

	// 2) xx服务...


	// 3. 启动服务器
	errc := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()

	// HTTP transport.
	go func() {
		logger.Log("transport", "HTTP", "addr", httpAddr)
		errc <- http.ListenAndServe(httpAddr, r)
	}()

	// Run!
	logger.Log("exit", <-errc)
}

func svcFactory(ctx context.Context, method, path string) sd.Factory {
	return func(instance string) (endpoint.Endpoint, io.Closer, error) {
		if !strings.HasPrefix(instance, "http") {
			instance = "http://" + instance
		}
		tgt, err := url.Parse(instance)
		fmt.Println("svcFactory url: ", tgt)
		if err != nil {
			return nil, nil, err
		}
		tgt.Path = path

		var (
			enc httptransport.EncodeRequestFunc
			dec httptransport.DecodeResponseFunc
		)
		enc, dec = encodeGetRequest, decodeGetResponse

		return httptransport.NewClient(method, tgt, enc, dec).Endpoint(), nil, nil
	}
}

func encodeGetRequest(_ context.Context, req *http.Request, request interface{}) error {

	data := request.(struct {
		UID string `json:"s"`
	})
	req.URL.Path = strings.Replace(req.URL.Path, "{UID}", data.UID, -1)

	//json
	//var buf bytes.Buffer
	//if err := json.NewEncoder(&buf).Encode(request); err != nil {
	//	return err
	//}
	//req.Body = ioutil.NopCloser(&buf)

	return nil
}

func encodeJSONResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func decodeGetResponse(ctx context.Context, resp *http.Response) (interface{}, error) {
	var getuserResponse struct {
		Name   string `json:"v"`
		Err error `json:"err,omitempty"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&getuserResponse); err != nil {
		return nil, err
	}
	return getuserResponse, nil
}

func decodeGetRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	var getuserRequest struct {
		UID string `json:"s"`
	}

	vars := mux.Vars(req)
	id, ok := vars["UID"]

	if !ok {
		return nil, errBadRoute
	}
	getuserRequest.UID = id
	fmt.Println("request: ", id)
	return getuserRequest, nil
}

var errBadRoute = errors.New("error argument")