package service

import (
	"context"
	"fmt"
	prod "github.com/TyrellJing/Hermes/go-kit-service/gorpc/server/proto"
	"io"
	"time"
)

type UserService struct {

}

func (u *UserService) GetUserScore(ctx context.Context, r *prod.UserScoreRequest) (*prod.UserScoreRepsonse, error) {
	var score int32 = 101
	users := make ([]*prod.UserInfo, 0)
	for _, user := range r.Users {
		user.UserScore = score
		score++
		users = append(users, user)
	}
	return &prod.UserScoreRepsonse{
		Users: users,
	}, nil
}

func (u *UserService) GetUserScoreByServerStream(r *prod.UserScoreRequest, stream prod.UserService_GetUserScoreByServerStreamServer) error {
	var score int32 = 101
	users := make ([]*prod.UserInfo, 0)
	for index, user := range r.Users {
		user.UserScore = score
		score++
		users = append(users, user)
		if (index+1) % 2 == 0 && index > 0 {
			err := stream.Send(&prod.UserScoreRepsonse{
				Users: users,
			})
			if err != nil {
				return err
			}
			users = (users)[0:0]
		}
		time.Sleep(time.Second * 1)
	}
	if len(users) > 0 {
		err := stream.Send(&prod.UserScoreRepsonse{
			Users: users,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *UserService) GetUserScoreByClientStream(stream prod.UserService_GetUserScoreByClientStreamServer) error {
	var score int32 = 101
	users := make ([]*prod.UserInfo, 0)
	for  {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&prod.UserScoreRepsonse{
				Users: users,
			})
		}
		if err != nil {
			return err
		}
		for _, user := range req.Users {
			user.UserScore = score
			score++
			users = append(users, user)
		}
	}
	return nil
}

func (u *UserService) GetUserScoreByTWStream(stream prod.UserService_GetUserScoreByTWStreamServer) error {
	var score int32 = 101
	users := make ([]*prod.UserInfo, 0)
	for  {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		for _, user := range req.Users {
			user.UserScore = score
			score++
			users = append(users, user)
		}
		err = stream.Send(&prod.UserScoreRepsonse{
			Users: users,
		})
		if err != nil {
			fmt.Println(err.Error())
		}
		users = (users)[0:0]
	}
	return nil
}