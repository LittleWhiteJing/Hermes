package Services

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

const secKey = "user_123"

type UserClaim struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type IAccessService interface {
	GetAccessToken(username string, password string) (string, error)
}

type AccessService struct {

}

func (as AccessService) GetAccessToken (username string, password string) (string, error)  {
	if username == "Yutaka" && password == "admin" {
		userinfo := &UserClaim { Username: username }
		userinfo.ExpiresAt = time.Now().Add(time.Second * 1000).Unix()
		oToken := jwt.NewWithClaims(jwt.SigningMethodHS256, userinfo)
		token, err := oToken.SignedString([]byte(secKey))
		return token, err
	}
	return "", fmt.Errorf("Error Username or Password")
}
