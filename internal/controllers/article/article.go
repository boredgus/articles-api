package article

import (
	"a-article/internal/controllers"
	"a-article/internal/models"
)

type ArticleController interface {
	Create(ctx controllers.Context) error
	Get(ctx controllers.Context) error
	GetForUser(ctx controllers.Context) error
	Update(ctx controllers.Context) error
	Delete(ctx controllers.Context) error
	UpdateReactionForArticle(ctx controllers.Context) error
}

func NewArticleController(article models.ArticleModel) ArticleController {
	return Article{articleModel: article}
}

type Article struct {
	articleModel models.ArticleModel
}
