package Services

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"errors"
)

func GetUserInfoRequest(ctx context.Context, request *http.Request, r interface{}) error {
	userRequest := r.(UserRequest)
	request.URL.Path += "/user/" + strconv.Itoa(userRequest.Uid)
	return nil
}

func GetUserInfoResponse(ctx context.Context, response *http.Response) (r interface{}, err error) {
	if response.StatusCode > 400 {
		return nil, errors.New("no data")
	}
	var userResponse UserResponse
	err = json.NewDecoder(response.Body).Decode(&userResponse)
	if err != nil {
		return nil, err
	}
	return userResponse, nil
}