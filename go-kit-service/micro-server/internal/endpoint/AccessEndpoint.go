package endpoint

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"micro-server/internal/service"
	"micro-server/util"
)

type AccessRequest struct {
	Username string
	Password string
	Method   string
}

type AccessResponse struct {
	Status int
	Token  string
}

func AccessEndpoint (accessService service.IAccessService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(AccessRequest)
		result := AccessResponse{ Status: 200 }
		if r.Method == "POST" {
			token, err := accessService.GetAccessToken(r.Username, r.Password)
			if err != nil {
				return nil, err
			} else {
				result.Token = token
				return result, nil
			}
		}
		return nil, util.NewAppError(403, "Forbidden")
	}
}