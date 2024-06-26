package article

import (
	"a-article/internal/controllers"
	"a-article/internal/models"
	"a-article/internal/views"
	"errors"
	"fmt"
	"net/http"
)

// swagger:parameters get_article
// nolint:unused
type getArticleParams struct {
	// article identifier
	// in: path
	// required: true
	ArticleId string `json:"article_id"`
}

// success
// swagger:response getArticleResp200
// nolint:unused
type getArticleResp200 struct {
	// in: body
	body views.Article
}

// swagger:route GET /articles/{article_id} articles get_article
// gets article
// ---
// Gets article by id.
// responses:
//
//	200: getArticleResp200
//	404: articleNotFound404
//	500: commonError
func (a Article) Get(ctx controllers.Context) error {
	article, err := a.articleModel.Get(ctx.PathParam("article_id"))
	if errors.Is(err, models.NotFoundErr) {
		e := ctx.JSON(http.StatusNotFound, controllers.ErrorBody{Error: err.Error()})
		return fmt.Errorf("%v: %w", e, err)
	}
	if err != nil {
		e := ctx.NoContent(http.StatusInternalServerError)
		return fmt.Errorf("%v: %w", e, err)
	}

	return ctx.JSON(http.StatusOK, views.NewArticleView(article))
}
