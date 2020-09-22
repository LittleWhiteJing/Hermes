package endpoint

import (
	"context"
	"fmt"
	"micro-server/internal/service"
	"micro-server/util"
	"net/http"
	"os"
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

func GenUserEndpoint(userService service.IUserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(UserRequest)
		fmt.Println("Current User: ", ctx.Value("Username"))
		result := "nothing"
		if r.Method == "GET" {
			result = userService.GetUsername(r.Uid) + strconv.Itoa(os.Getpid())
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



