package internal

import (
	"a-article/internal/controllers/article"
	user "a-article/internal/controllers/user"
	"a-article/internal/gateways"
	"a-article/internal/models"
	"a-article/pkg/msgbroker"
)

type AppController struct {
	User    user.UserController
	Article article.ArticleController
}

type AppParams struct {
	MainStore     gateways.Store
	StatsStore    gateways.Store
	CacheStore    gateways.CacheStore
	MessageBroker msgbroker.Broker
}

func NewAppController(params AppParams) AppController {
	return AppController{
		User: user.NewUserController(models.NewUserModel(gateways.NewUserRepository(params.MainStore), params.MessageBroker)),
		Article: article.NewArticleController(
			models.NewArticleModel(
				gateways.NewCachedArticleRepository(
					gateways.NewArticleRepository(params.MainStore, params.StatsStore),
					params.CacheStore),
			)),
	}
}
