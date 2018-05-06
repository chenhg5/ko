package gateway

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func RunServer(logger log.Logger, httpAddr string, router *Router) {
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
		errc <- http.ListenAndServe(httpAddr, router.r)
	}()

	// Run!
	logger.Log("exit", <-errc)
}
