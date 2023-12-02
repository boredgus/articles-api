package article

import (
	"errors"
	"fmt"
	"net/http"
	"user-management/internal/controllers"
	"user-management/internal/domain"
	"user-management/internal/models"
	"user-management/internal/views"
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

// user_oid or password is invalid
// swagger:response createArticleResp401
// nolint:unused
type createArticleResp401 struct {
	// in: body
	body controllers.ErrorBody
}

// invalid article data provided
// swagger:response createArticleResp400
// nolint:unused
type authResp401 struct {
	// in: body
	body controllers.ErrorBody
}

// swagger:route POST /articles articles create_article
// creates new article
// ---
// Checks user with provided `user_oid` and `password` exists and validates article data.
// Then saves and returns created item.
//
// security:
//
//   - userOId:
//     password:
//
// responses:
//
//	201: createArticleResp201
//	401: createArticleResp401
//	400: createArticleResp400
//	500: commonError
func (a Article) Create(ctx controllers.Context) error {
	userOId := ctx.Request().Header.Get("User-OId")
	err := a.userModel.Exists(userOId, ctx.QueryParams().Get("password"))
	if errors.Is(err, models.UserNotFoundErr) || errors.Is(err, models.InvalidAuthParameterErr) {
		e := ctx.JSON(http.StatusUnauthorized, controllers.ErrorBody{Error: "invalid user_oid or password"})
		return fmt.Errorf("%v: %w", e, err)
	}
	if err != nil {
		e := ctx.NoContent(http.StatusInternalServerError)
		return fmt.Errorf("%v: %w", e, err)
	}

	var data ArticleData
	err = ctx.Bind(&data)
	if err != nil {
		e := ctx.JSON(http.StatusBadRequest, controllers.ErrorBody{Error: "failed to parse article"})
		return fmt.Errorf("%v: %w", e, err)
	}
	article := domain.Article{Theme: data.Theme, Text: data.Text, Tags: data.Tags}
	if len(article.Tags) == 0 {
		article.Tags = []string{}
	}
	err = a.articleModel.Create(userOId, &article)
	if errors.Is(err, models.InvalidArticleErr) {
		e := ctx.JSON(http.StatusBadRequest, controllers.ErrorBody{Error: err.Error()})
		return fmt.Errorf("%v: %w", e, err)
	}
	if err != nil {
		e := ctx.NoContent(http.StatusInternalServerError)
		return fmt.Errorf("%v: %w", e, err)
	}
	return ctx.JSON(http.StatusCreated, views.NewArticleView(article))
}
