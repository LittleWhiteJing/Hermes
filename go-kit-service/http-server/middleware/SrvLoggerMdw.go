package middleware

import (
	"context"
	userendpoint "micro-server/internal/endpoint"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
)

func SrvLogger(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			r := request.(userendpoint.UserRequest)
			logger.Log("method", r.Method, "event", "get_user", "userid", r.Uid)
			return next(ctx, request)
		}
	}
}
