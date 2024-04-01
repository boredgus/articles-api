package user

import (
	"a-article/internal/controllers"
	"a-article/internal/models"
)

type UserController interface {
	Register(ctx controllers.Context) error
	Authorize(ctx controllers.Context) error
	Delete(ctx controllers.Context) error
	UpdateRole(ctx controllers.Context) error
}

func NewUserController(userModel models.UserModel) UserController {
	return User{userModel: userModel}
}

type User struct {
	userModel models.UserModel
}
