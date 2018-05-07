package middleware

import (
	"fmt"
	"github.com/go-kit/kit/metrics"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"ko/services/ucenter"
	"time"
)

var (
	fieldKeys    = []string{"method", "error"}
	requestCount = kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "my_group",
		Subsystem: "ucenter_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)
	requestLatency = kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "my_group",
		Subsystem: "ucenter_service",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)
	countResult = kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "my_group",
		Subsystem: "ucenter_service",
		Name:      "count_result",
		Help:      "The result of each count method.",
	}, []string{}) // no fields here
)

func InstrumentingMiddleware() ucenter.ServiceMiddleware {
	return func(next ucenter.UcenterServiceInterface) ucenter.UcenterServiceInterface {
		return instrumentingMiddleware{requestCount, requestLatency, countResult, next}
	}
}

type instrumentingMiddleware struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	countResult    metrics.Histogram
	ucenter.UcenterServiceInterface
}

func (mw instrumentingMiddleware) GetUser(s string) (output string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "uppercase", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	output, err = mw.UcenterServiceInterface.GetUser(s)
	return
}
