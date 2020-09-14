package service

import "errors"

type IUserService interface {
	GetUsername (uid int) string
	DelUserinfo (uid int) error
}

type UserService struct {

}

func (u UserService) GetUsername(uid int) string {
	if uid == 101 {
		return "Yutaka"
	}
	return "Guest"
}

func (u UserService) DelUserinfo(uid int) error {
	if uid == 101 {
		return errors.New("无权限")
	}
	return nil
}