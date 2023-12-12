package article

import (
	"errors"
	"fmt"
	"net/http"
	"user-management/internal/auth"
	"user-management/internal/controllers"
	"user-management/internal/domain"
	"user-management/internal/models"
	"user-management/internal/views"

	"github.com/golang-jwt/jwt/v5"
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

// success
// swagger:response updateArticleResp200
// nolint:unused
type updateArticleResp200 struct {
	// updated article
	// in: body
	Article domain.Article
}

// invalid article data
// swagger:response updateArticleResp400
// nolint:unused
type updateArticleResp400 struct {
	// in: body
	Body controllers.ErrorBody
}

// swagger:route PUT /articles/{article_id} articles update_article
// updates article
// ---
// - Checks whether article is owned by authorized user, validates provided article data, updates article and returns updates item.
//
// security:
//   - jwt:
//
// responses:
//
//	200: updateArticleResp200
//	400: updateArticleResp400
//	401: unauthorizedResp401
//	403: forbiddenResp403
//	404: notFoundResp404
//	500: commonError
func (a Article) Update(ctx controllers.Context) error {
	var data ArticleData
	err := ctx.Bind(&data)
	if err != nil {
		e := ctx.JSON(http.StatusBadRequest, controllers.ErrorBody{Error: "failed to parse article"})
		return fmt.Errorf("%v: %w", e, err)
	}
	article := domain.Article{
		OId:   ctx.PathParam("article_id"),
		Theme: data.Theme,
		Text:  data.Text,
		Tags:  data.Tags}
	if len(article.Tags) == 0 {
		article.Tags = []string{}
	}
	claims := ctx.Get("user").(*jwt.Token).Claims.(*auth.JWTClaims)
	err = a.articleModel.Update(claims.UserOId, claims.Role, &article)
	if errors.Is(err, models.ArticleNotFoundErr) {
		e := ctx.JSON(http.StatusNotFound, controllers.ErrorBody{Error: err.Error()})
		return fmt.Errorf("%v: %w", e, err)
	}
	if errors.Is(err, models.NotEnoughRightsErr) {
		e := ctx.JSON(http.StatusForbidden, controllers.ErrorBody{Error: err.Error()})
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
	return ctx.JSON(http.StatusOK, views.NewArticleView(article))
}
