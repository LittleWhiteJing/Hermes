package Services

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/endpoint"
)

type UserRequest struct {
	Uid int `json:"uid"`
	Method string
}

type UserResponse struct {
	Result string `json:"result"`
}

func GenUserEndpoint(userService IUserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(UserRequest)
		result := "nothing"
		if r.Method == "GET" {
			result = userService.GetUsername(r.Uid)
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