package Services

import (
	"context"
	"fmt"
	"golang.org/x/time/rate"
	"micro-server/util"
	"net/http"
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
				return nil, util.NewAppError(429,"To Many Requests")
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

func AppErrorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	contentType, body := "text/plain; charset=utf-8", []byte(err.Error())
	w.Header().Set("Content-Type", contentType)
	if appErr, ok := err.(*util.AppError); ok {
		w.WriteHeader(appErr.Code)
	} else {
		w.WriteHeader(500)
	}
	w.Write(body)
}



