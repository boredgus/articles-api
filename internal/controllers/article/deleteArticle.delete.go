package article

import (
	"errors"
	"fmt"
	"net/http"
	"user-management/internal/auth"
	"user-management/internal/controllers"
	"user-management/internal/models"

	"github.com/golang-jwt/jwt/v5"
)

// user deleted
// swagger:response deleteUserResp200
// nolint:unused
type deleteUserResp200 struct{}

// swagger:parameters delete_article
// nolint:unused
type deleteArticleParams struct {
	// article identifier
	// in: path
	// required: true
	ArticleId string `json:"article_id"`
}

// swagger:route DELETE /articles/{article_id} articles delete_article
// deletes article
// ---
// - Checks whether article is owned by authorized user, validates provided article data, updates article and returns updates item.
//
// security:
//   - jwt:
//
// responses:
//
//	200: deleteUserResp200
//	401: unauthorizedResp401
//	403: forbiddenResp403
//	404: notFoundResp404
//	500: commonError
func (a Article) Delete(ctx controllers.Context) error {
	claims := ctx.Get("user").(*jwt.Token).Claims.(*auth.JWTClaims)
	err := a.articleModel.Delete(claims.UserOId, claims.Role, ctx.PathParam("article_id"))
	if errors.Is(err, models.ArticleNotFoundErr) {
		e := ctx.JSON(http.StatusNotFound, controllers.ErrorBody{Error: err.Error()})
		return fmt.Errorf("%v: %w", e, err)
	}
	if errors.Is(err, models.NotEnoughRightsErr) {
		e := ctx.JSON(http.StatusForbidden, controllers.ErrorBody{Error: err.Error()})
		return fmt.Errorf("%v: %w", e, err)
	}
	if err != nil {
		e := ctx.NoContent(http.StatusInternalServerError)
		return fmt.Errorf("%v: %w", e, err)
	}
	return ctx.NoContent(http.StatusOK)
}