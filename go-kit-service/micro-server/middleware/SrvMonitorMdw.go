package middleware

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/metrics"
	"time"
)

func SrvMonitor(requestCount metrics.Counter, requestLatency metrics.Histogram) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			begin := time.Now()
			lvs := []string{"method", "Add"}
			requestCount.With(lvs...).Add(1)
			requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
			return next(ctx, request)
		}
	}
}
