package article

import (
	"user-management/internal/controllers"
	"user-management/internal/models"
)

type ArticleController interface {
	Create(ctx controllers.Context) error
	GetForUser(ctx controllers.Context) error
	Update(ctx controllers.Context) error
}

func NewArticleController(user models.UserModel, article models.ArticleModel) ArticleController {
	return Article{userModel: user, articleModel: article}
}

type Article struct {
	articleModel models.ArticleModel
	userModel    models.UserModel
}
