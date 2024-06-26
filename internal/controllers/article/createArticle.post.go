package article

import (
	"a-article/internal/auth"
	"a-article/internal/controllers"
	"a-article/internal/domain"
	"a-article/internal/models"
	"a-article/internal/views"
	"errors"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

type ArticleData struct {
	// theme of article
	// required: true
	Theme string `json:"theme" form:"theme"`
	// content of article
	Text string `json:"text" form:"text"`
	// topics of article
	Tags []string `json:"tags" form:"tags"`
}

// swagger:parameters create_article
// nolint:unused
type createParameters struct {
	// article data
	// in: body
	// required: true
	Article ArticleData `json:"article"`
}

// successfully created
// swagger:response createArticleResp201
// nolint:unused
type authResp200 struct {
	// in: body
	body domain.Article
}

// swagger:route POST /articles articles create_article
// creates new article
// ---
// Checks user with provided `user_oid` and `password` exists and validates article data.
// Then creates and returns created item.
//
// security:
//   - jwt:
//
// responses:
//
//	201: createArticleResp201
//	401: unauthorizedResp401
//	400: invalidData400
//	500: commonError
func (a Article) Create(ctx controllers.Context) error {
	var data ArticleData
	err := ctx.Bind(&data)
	if err != nil {
		e := ctx.JSON(http.StatusBadRequest, controllers.ErrorBody{Error: "failed to parse article"})
		return fmt.Errorf("%v: %w", e, err)
	}
	article := domain.Article{Theme: data.Theme, Text: data.Text, Tags: data.Tags}
	if len(article.Tags) == 0 {
		article.Tags = []string{}
	}
	err = a.articleModel.Create(ctx.Get("user").(*jwt.Token).Claims.(*auth.JWTClaims).UserOId, &article)
	if errors.Is(err, models.InvalidDataErr) {
		e := ctx.JSON(http.StatusBadRequest, controllers.ErrorBody{Error: err.Error()})
		return fmt.Errorf("%v: %w", e, err)
	}
	if err != nil {
		e := ctx.NoContent(http.StatusInternalServerError)
		return fmt.Errorf("%v: %w", e, err)
	}
	return ctx.JSON(http.StatusCreated, views.NewArticleView(article))
}
