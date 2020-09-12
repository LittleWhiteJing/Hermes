package Services

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func DecodeUserRequest (ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	if userid, ok := vars["userid"]; ok {
		uid, _ := strconv.Atoi(userid)
		return UserRequest{
			Uid: uid,
			Method: r.Method,
		}, nil
	}
	return nil, errors.New("参数错误！")
}

func EncodeUserResponse (ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Context-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}
