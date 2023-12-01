package article

import (
	"errors"
	"fmt"
	"net/http"
	"user-management/internal/controllers"
	"user-management/internal/domain"
	"user-management/internal/models"
)

// swagger:parameters update_article
// nolint:unused
type updateArticleParams struct {
	// article identifier
	// in: path
	// required: true
	ArticleId string `json:"article_id"`
	// article data to update
	// in: body
	// required: true
	Article ArticleData `json:"article"`
}

// swagger:response updateArticleResp200
// nolint:unused
type updateArticleResp200 struct {
	// updated article
	// in: body
	// required: true
	Article domain.Article
}

// swagger:route PUT /articles/{article_id} articles update_article
// updates article
// ---
// - Checks whether article is owned by authorized user, validates provided article data, updates article and returns updates item.
//
// security:
//   - BasicAuth:
//
// responses:
//
//	200: updateArticleResp200
func (a Article) Update(ctx controllers.Context) error {
	var data ArticleData
	err := ctx.Bind(&data)
	if err != nil {
		e := ctx.JSON(http.StatusBadRequest, controllers.ErrorBody{Error: "failed to parse article"})
		return fmt.Errorf("%v: %w", e, err)
	}
	fmt.Printf("> data: %+v", data)
	article := domain.Article{
		OId:   ctx.PathParam("article_id"),
		Theme: data.Theme,
		Text:  data.Text,
		Tags:  data.Tags}
	if len(article.Tags) == 0 {
		article.Tags = []string{}
	}
	fmt.Printf("> article data: %+v\n", article)

	err = a.articleModel.Update(ctx.Request().Header.Get("Username"), &article)
	if errors.Is(err, models.UserIsNotAnOwnerErr) {
		e := ctx.JSON(http.StatusNotFound, controllers.ErrorBody{Error: err.Error()})
		return fmt.Errorf("%v: %w", e, err)
	}
	if errors.Is(err, models.InvalidArticleErr) {
		e := ctx.JSON(http.StatusBadRequest, controllers.ErrorBody{Error: err.Error()})
		return fmt.Errorf("%v: %w", e, err)
	}
	if err != nil {
		e := ctx.NoContent(http.StatusInternalServerError)
		return fmt.Errorf("%v: %w", e, err)
	}
	return ctx.JSON(http.StatusOK, article)
}
