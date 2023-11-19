package infrastructure

import (
	"user-management/internal/controllers"
	"user-management/internal/gateways"
	"user-management/internal/models"
)

func NewAppController(store gateways.Store) controllers.AppController {
	return controllers.AppController{
		User: controllers.NewLoginController(models.NewUserModel(gateways.NewUserRepository(store))),
	}
}
