package internal

import (
	"a-article/internal/controllers/article"
	user "a-article/internal/controllers/user"
	"a-article/internal/gateways"
	"a-article/internal/models"
)

type AppController struct {
	User    user.UserController
	Article article.ArticleController
}

func NewAppController(mainStore, statsStore gateways.Store, cacheStore gateways.CacheStore) AppController {
	return AppController{
		User: user.NewUserController(models.NewUserModel(gateways.NewUserRepository(mainStore))),
		Article: article.NewArticleController(
			models.NewArticleModel(
				gateways.NewCachedArticleRepository(
					gateways.NewArticleRepository(mainStore, statsStore),
					cacheStore),
			)),
	}
}
