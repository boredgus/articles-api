package internal

import (
	"user-management/internal/controllers/article"
	user "user-management/internal/controllers/login"
	"user-management/internal/gateways"
	"user-management/internal/models"
)

type AppController struct {
	Login   user.LoginController
	Article article.ArticleController
	User    models.UserModel
}

func NewAppController(store gateways.Store) AppController {
	userModel := models.NewUserModel(gateways.NewUserRepository(store))
	return AppController{
		User:  userModel,
		Login: user.NewLoginController(userModel),
		Article: article.NewArticleController(
			userModel,
			models.NewArticleModel(gateways.NewArticleRepository(store))),
	}
}
