package main

import (
	"net/http"
	"shop/services/ucenter"
	kitlog "github.com/go-kit/kit/log"
	"os"
	"os/signal"
	"syscall"
	"fmt"
	"flag"
)

func main() {

	var addr = "4001"
	var httpAddr = flag.String("http.addr", ":"+addr, "HTTP listen address")

	ucenterSvc := ucenter.UcenterService{}

	var logger kitlog.Logger
	logger = kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(os.Stderr))
	logger = kitlog.With(logger, "ts", kitlog.DefaultTimestampUTC)
	httpLogger := kitlog.With(logger, "component", "http")

	mux := http.NewServeMux()

	mux.Handle("/ucenter/v1/", ucenter.MakeHandler(ucenterSvc, httpLogger))

	http.Handle("/", accessControl(mux))

	errs := make(chan error, 2)
	go func() {
		logger.Log("transport", "http", "address", *httpAddr, "msg", "listening")
		errs <- http.ListenAndServe(*httpAddr, nil)
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