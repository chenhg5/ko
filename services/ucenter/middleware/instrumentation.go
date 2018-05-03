package middleware

import (
	"github.com/go-kit/kit/metrics"
	"fmt"
	"time"
	"ko/services/ucenter"
)


func InstrumentingMiddleware(
	requestCount metrics.Counter,
	requestLatency metrics.Histogram,
	countResult metrics.Histogram,
) ucenter.ServiceMiddleware {
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