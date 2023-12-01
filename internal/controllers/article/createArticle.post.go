package article

import (
	"errors"
	"fmt"
	"net/http"
	"user-management/internal/controllers"
	"user-management/internal/domain"
	"user-management/internal/models"
)

type ArticleData struct {
	// theme of article
	// required: true
	Theme string `json:"theme"`
	// content of article
	Text string `json:"text"`
	// topics of article
	Tags []string `json:"tags"`
}

// swagger:parameters create_article
// nolint:unused
type createParameters struct {
	// unique user identifier
	// in: header
	// required: true
	UserOId string `json:"user_oid"`
	// user password and article payload
	// in: body
	// required: true
	Body struct {
		// user password
		// required: true
		Password string `json:"password"`
		ArticleData
	} `json:"payload"`
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
// responses:
//
//	201: createArticleResp201
//	401: createArticleResp401
//	400: createArticleResp400
//	500: commonError
func (a Article) Create(ctx controllers.Context) error {
	formParams, err := ctx.FormParams()
	if err != nil {
		e := ctx.JSON(http.StatusBadRequest, controllers.ErrorBody{Error: "failed to parse form params"})
		return fmt.Errorf("%v: %w", e, err)
	}
	userOId := ctx.Request().Header.Get("user_oid")
	err = a.userModel.Exists(userOId, formParams.Get("password"))
	if errors.Is(err, models.UserNotFoundErr) || errors.Is(err, models.InvalidAuthParameterErr) {
		e := ctx.JSON(http.StatusUnauthorized, controllers.ErrorBody{Error: "invalid user_oid or password"})
		return fmt.Errorf("%v: %w", e, err)
	}
	if err != nil {
		e := ctx.NoContent(http.StatusInternalServerError)
		return fmt.Errorf("%v: %w", e, err)
	}

	article := domain.Article{
		Theme: formParams.Get("theme"),
		Text:  formParams.Get("text"),
		Tags:  formParams["tags"]}

	err = a.articleModel.Create(userOId, &article)
	if errors.Is(err, models.InvalidArticleErr) {
		e := ctx.JSON(http.StatusBadRequest, controllers.ErrorBody{Error: err.Error()})
		return fmt.Errorf("%v: %w", e, err)
	}
	if err != nil {
		e := ctx.NoContent(http.StatusInternalServerError)
		return fmt.Errorf("%v: %w", e, err)
	}
	return ctx.JSON(http.StatusCreated, article)
}
