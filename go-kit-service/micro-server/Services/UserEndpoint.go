package Services

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/log"
	"golang.org/x/time/rate"
	"micro-server/util"
	"net/http"
	"strconv"

	"github.com/go-kit/kit/endpoint"
)

type UserRequest struct {
	Uid int `json:"uid"`
	Method string
	Token  string
}

type UserResponse struct {
	Result string `json:"result"`
}

func UserAuth() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			r := request.(UserRequest)
			if r.Token != "" {
				uc := UserClaim{}
				getToken, err := jwt.ParseWithClaims(r.Token, &uc, func(token *jwt.Token) (interface{}, error) {
					return []byte(secKey), nil
				})
				if err != nil {
					return nil, util.NewAppError(403, "Invalid Token")
				}
				if getToken != nil && getToken.Valid {
					newCtx := context.WithValue(ctx, "Username", getToken.Claims.(*UserClaim).Username)
					return next(newCtx, request)
				}
			}
			return nil, util.NewAppError(403, "Forbidden")
		}
	}
}

func SrvLogger(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			r := request.(UserRequest)
			logger.Log("method", r.Method, "event", "get_user", "userid", r.Uid)
			return next(ctx, request)
		}
	}
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
		fmt.Println("Current User: ", ctx.Value("Username"))
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



