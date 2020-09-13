package Services

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/time/rate"
	"micro-server/util"
	"strconv"

	"github.com/go-kit/kit/endpoint"
)

type UserRequest struct {
	Uid int `json:"uid"`
	Method string
}

type UserResponse struct {
	Result string `json:"result"`
}

func RateLimit(limit *rate.Limiter) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if !limit.Allow() {
				return nil, errors.New("To Many Requests")
			}
			return next(ctx, request)
		}
	}
}

func GenUserEndpoint(userService IUserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(UserRequest)
		result := "nothing"
		if r.Method == "GET" {
			result = userService.GetUsername(r.Uid) + strconv.Itoa(util.ServicePort)
		} else if r.Method == "DELETE" {
			err := userService.DelUserinfo(r.Uid)
			if err != nil {
				result = err.Error()
			} else {
				result = fmt.Sprintf("用户id:%d，删除成功", r.Uid)
			}
		}

		return UserResponse{ Result: result }, nil
	}
}