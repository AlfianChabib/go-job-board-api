package service

import "AlfianChabib/go-job-board-api/internal/model/web"

type AuthService interface {
	Register(data web.RegisterRequest) error
	Login(data web.LoginRrequest) error
	Logout() error
}
