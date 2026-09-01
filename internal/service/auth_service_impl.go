package service

import (
	"AlfianChabib/go-job-board-api/internal/model/web"
)

type AuthServiceImpl struct{}

func NewAuthService() *AuthServiceImpl {
	return &AuthServiceImpl{}
}

func (a *AuthServiceImpl) Register(data web.RegisterRequest) error {
	panic("TODO: Implement")
}

func (a *AuthServiceImpl) Login(data web.LoginRrequest) error {
	panic("TODO: Implement")
}

func (a *AuthServiceImpl) Logout() error {
	panic("TODO: Implement")
}
