package ucenter

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"os"
	"os/signal"
	"fmt"
	kitlog "github.com/go-kit/kit/log"
	"log"
	"net/http"
	"time"
	"context"
	"io/ioutil"
)

func RunServer(mux *http.ServeMux, logger kitlog.Logger, httpAddr string)  {
	http.Handle("/", accessControl(mux))
	http.Handle("/metrics", promhttp.Handler())

	srv := &http.Server{
		Addr:    httpAddr,
		Handler: mux,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")

	pid := fmt.Sprintf("%d", os.Getpid())
	_, openErr := os.OpenFile((*GetConfig())["pid_path"], os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if openErr == nil {
		ioutil.WriteFile((*GetConfig())["pid_path"], []byte(pid), 0)
	}
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