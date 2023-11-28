package internal

import (
	user "user-management/internal/controllers/login"
	"user-management/internal/gateways"
	"user-management/internal/models"
)

type AppController struct {
	User user.LoginController
}

func NewAppController(store gateways.Store) AppController {
	return AppController{
		User: user.NewLoginController(models.NewUserModel(gateways.NewUserRepository(store))),
	}
}
