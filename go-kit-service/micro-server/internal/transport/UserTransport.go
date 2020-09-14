package transport

import (
	"context"
	"encoding/json"
	"errors"
	"micro-server/internal/endpoint"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func DecodeUserRequest (ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	var userRequest endpoint.UserRequest
	if userid, ok := vars["userid"]; ok {
		uid, _ := strconv.Atoi(userid)
		userRequest.Uid = uid
	} else {
		return nil, errors.New("Invalid Params")
	}
	if token := r.URL.Query().Get("token"); token != "" {
		userRequest.Token = token
	} else {
		return nil, errors.New("Invalid Params")
	}
	return userRequest, nil
}

func EncodeUserResponse (ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Context-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}
