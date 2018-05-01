package main

import (
	"github.com/go-kit/kit/sd/etcdv3"
	"time"
	"github.com/go-kit/kit/log"
	"context"
	"os"
	"os/signal"
	"syscall"
	"fmt"
	"net/http"
	kitlog "github.com/go-kit/kit/log"
	"ko/services/ucenter"
)

func main() {
	// 1. 配置
	var (
		etcdServer = "localhost:2379"      // in the change from v2 to v3, the schema is no longer necessary if connecting directly to an etcd v3 instance
		ctx        = context.Background()
		prefix     = "/svc/ucenter/"  // known at compile time
		instance   = "localhost:4002"       // taken from runtime or platform, somehow
		key        = prefix + instance    // should be globally unique
		value      = "http://" + instance // based on our transport
		httpAddr   = ":4002"
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

	// Build the registrar.
	registrar := etcdv3.NewRegistrar(etcdClient, etcdv3.Service{
		Key:   key,
		Value: value,
	}, log.NewNopLogger())

	// Register our instance.
	registrar.Register()

	// At the end of our service lifecycle, for example at the end of func main,
	// we should make sure to deregister ourselves. This is important! Don't
	// accidentally skip this step by invoking a log.Fatal or os.Exit in the
	// interim, which bypasses the defer stack.
	defer registrar.Deregister()

	ucenterSvc := ucenter.UcenterService{}

	var logger kitlog.Logger
	logger = kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(os.Stderr))
	logger = kitlog.With(logger, "ts", kitlog.DefaultTimestampUTC)
	httpLogger := kitlog.With(logger, "component", "http")

	mux := http.NewServeMux()

	mux.Handle("/svc/ucenter/v1/", ucenter.MakeHandler(ucenterSvc, httpLogger))

	http.Handle("/", accessControl(mux))

	errs := make(chan error, 2)
	go func() {
		logger.Log("transport", "http", "address", httpAddr, "msg", "listening")
		errs <- http.ListenAndServe(httpAddr, nil)
	}()
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	logger.Log("terminated", <-errs)
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}