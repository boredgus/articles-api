package internal

import (
	"user-management/internal/controllers/article"
	user "user-management/internal/controllers/user"
	"user-management/internal/gateways"
	"user-management/internal/models"
)

type AppController struct {
	User    user.UserController
	Article article.ArticleController
}

func NewAppController(store gateways.Store) AppController {
	userModel := models.NewUserModel(gateways.NewUserRepository(store))
	return AppController{
		User: user.NewUserController(userModel),
		Article: article.NewArticleController(
			userModel,
			models.NewArticleModel(gateways.NewArticleRepository(store))),
	}
}
