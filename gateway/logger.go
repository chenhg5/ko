package gateway

import (
	"os"
	"io"
	"github.com/go-kit/kit/log"
)

func InitLogger() log.Logger {
	Accessfile, err := os.OpenFile("access.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer Accessfile.Close()

	var kitlogger log.Logger
	kitlogger = log.NewLogfmtLogger(io.MultiWriter(os.Stderr, Accessfile))
	kitlogger = log.With(kitlogger, "ts", log.DefaultTimestampUTC)

	return log.With(kitlogger, "component", "http")
}