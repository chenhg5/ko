package main

import (
	"time"
	"github.com/go-kit/kit/sd/etcdv3"
	"context"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/lb"
	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	httptransport "github.com/go-kit/kit/transport/http"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"fmt"
	"ko/gateway"
	"io"
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

	// 2. 日志系统

	Accessfile, err := os.OpenFile("access.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer Accessfile.Close()

	var kitlogger log.Logger
	kitlogger = log.NewLogfmtLogger(io.MultiWriter(os.Stderr, Accessfile))
	kitlogger = log.With(kitlogger, "ts", log.DefaultTimestampUTC)
	logger := log.With(kitlogger, "component", "http")

	// 3. 服务发现
	etcdClient, err := etcdv3.NewClient(ctx, []string{etcdServer}, options)
	if err != nil {
		panic(err)
	}

	// 1) 用户中心服务

	// 创造实例
	ucenterPrefix := "/svc/ucenter"
	ucenterInstancer, err := etcdv3.NewInstancer(etcdClient, ucenterPrefix, logger)
	if err != nil {
		panic(err)
	}

	// 路由控制器构造
	ucenterfactory := gateway.SvcFactory(ctx, "GET", "/svc/ucenter/v1/user/{param}")
	ucenterendpointer := sd.NewEndpointer(ucenterInstancer, ucenterfactory, logger)
	ucenterbalancer := lb.NewRoundRobin(ucenterendpointer)
	ucenterretry := lb.Retry(3, 3*time.Second, ucenterbalancer)

	// 路由
	r := mux.NewRouter()
	r.Handle("/svc/ucenter/v1/user/{param}", httptransport.NewServer(ucenterretry, gateway.DecodeGetRequest, gateway.EncodeJSONResponse))

	// 2) 订单服务...
	orderPrefix := "/svc/order"
	orderInstancer, err := etcdv3.NewInstancer(etcdClient, orderPrefix, logger)
	if err != nil {
		panic(err)
	}

	// 路由控制器构造
	orderfactory := gateway.SvcFactory(ctx, "POST", "/svc/order/v1/order")
	orderendpointer := sd.NewEndpointer(orderInstancer, orderfactory, logger)
	orderbalancer := lb.NewRoundRobin(orderendpointer)
	orderretry := lb.Retry(3, 3*time.Second, orderbalancer)

	// 路由
	r.Handle("/svc/order/v1/order", httptransport.NewServer(orderretry, gateway.DecodeJsonRequest, gateway.EncodeJSONResponse))

	// 3) xx服务...


	// 4. 启动服务器
	errc := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()

	// pid := fmt.Sprintf("%d", os.Getpid())
	// // TODO: pid文件位置放在全局设置中
	// _, openErr := os.OpenFile("./gateway/pid", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	// if openErr == nil {
	// 	ioutil.WriteFile("./gateway/pid", []byte(pid), 0)
	// }

	// HTTP transport.
	go func() {
		logger.Log("transport", "HTTP", "addr", httpAddr)
		errc <- http.ListenAndServe(httpAddr, r)
	}()

	// Run!
	logger.Log("exit", <-errc)
}