package middleware

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"golang.org/x/time/rate"
	"micro-server/util"
)

func RateLimit(limit *rate.Limiter) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if !limit.Allow() {
				return nil, util.NewAppError(429,"To Many Requests")
			}
			return next(ctx, request)
		}
	}
}
