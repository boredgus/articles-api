package user

import (
	"user-management/internal/controllers"
	"user-management/internal/models"
)

type LoginController interface {
	Register(ctx controllers.Context) error
	Authorize(ctx controllers.Context) error
}

func NewLoginController(userModel models.UserModel) LoginController {
	return Login{userModel: userModel}
}

type Login struct {
	userModel models.UserModel
}
