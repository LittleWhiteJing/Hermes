package middleware

import (
	"context"
	userendpoint "micro-server/internal/endpoint"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/endpoint"
	"micro-server/internal/service"
	"micro-server/util"
)

func JwtAuth() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			r := request.(userendpoint.UserRequest)
			if r.Token != "" {
				uc := service.UserClaim{}
				getToken, err := jwt.ParseWithClaims(r.Token, &uc, func(token *jwt.Token) (interface{}, error) {
					return []byte(service.SecKey), nil
				})
				if err != nil {
					return nil, util.NewAppError(403, "Invalid Token")
				}
				if getToken != nil && getToken.Valid {
					newCtx := context.WithValue(ctx, "Username", getToken.Claims.(*service.UserClaim).Username)
					return next(newCtx, request)
				}
			}
			return nil, util.NewAppError(403, "Forbidden")
		}
	}
}
