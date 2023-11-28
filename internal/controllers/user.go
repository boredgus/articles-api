package controllers

import (
	"user-management/internal/models"
)

type LoginController interface {
	Register(ctx Context) error
	Authorize(ctx Context) error
}

func NewLoginController(userModel models.UserModel) LoginController {
	return Login{userModel: userModel}
}

type Login struct {
	userModel models.UserModel
}
