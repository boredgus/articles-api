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

func NewAppController(mainStore, statsStore gateways.Store) AppController {
	return AppController{
		User: user.NewUserController(models.NewUserModel(gateways.NewUserRepository(mainStore))),
		Article: article.NewArticleController(
			models.NewArticleModel(gateways.NewArticleRepository(mainStore, statsStore))),
	}
}
