package article

import (
	"a-article/internal/auth"
	"a-article/internal/controllers"
	"a-article/internal/models"
	"errors"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

type reactionPayload struct {
	// emoji to express reaction for article. Empty string stands for withdrawing of reaction.
	// example: 😂
	Emoji string `json:"emoji" form:"emoji"`
}

// swagger:parameters reaction_for_article
// nolint:unused
type rateArticleParams struct {
	// article identifier
	// in: path
	// required: true
	ArticleId string `json:"article_id"`
	// reaction value to react on article
	// in: body
	// required: true
	Body reactionPayload `json:"payload"`
}

// swagger:route PUT /articles/{article_id}/reaction articles reaction_for_article
// updates reaction for article
// ---
// User is forbidden to react on own articles.
// security:
//   - jwt:
//
// responses:
//
//	200: respWithMessage
//	400: invalidData400
//	401: unauthorizedResp401
//	403: forbiddenResp403
//	404: articleNotFound404
//	500: commonError
func (a Article) UpdateReactionForArticle(ctx controllers.Context) error {
	var payload reactionPayload
	err := ctx.Bind(&payload)
	if err != nil {
		return fmt.Errorf("%v: %w", ctx.JSON(http.StatusBadRequest, controllers.ErrorBody{Error: "failed to bind data"}), err)
	}
	err = a.articleModel.UpdateReaction(
		ctx.Get("user").(*jwt.Token).Claims.(*auth.JWTClaims).UserOId,
		ctx.PathParam("article_id"),
		payload.Emoji,
	)
	if errors.Is(err, models.NotFoundErr) {
		e := ctx.JSON(http.StatusNotFound, controllers.ErrorBody{Error: err.Error()})
		return fmt.Errorf("%v: %w", e, err)
	}
	if errors.Is(err, models.NotEnoughRightsErr) {
		e := ctx.JSON(http.StatusForbidden, controllers.ErrorBody{Error: err.Error()})
		return fmt.Errorf("%v: %w", e, err)
	}
	if errors.Is(err, models.InvalidDataErr) {
		e := ctx.JSON(http.StatusBadRequest, controllers.ErrorBody{Error: err.Error()})
		return fmt.Errorf("%v: %w", e, err)
	}
	if err != nil {
		return fmt.Errorf("%v: %w", ctx.NoContent(http.StatusInternalServerError), err)
	}
	return ctx.JSON(http.StatusOK, controllers.InfoResponse{Message: fmt.Sprintf("reaction set to '%v'", payload.Emoji)})
}
