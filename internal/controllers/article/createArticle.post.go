package article

import (
	"errors"
	"fmt"
	"net/http"
	"user-management/internal/controllers"
	"user-management/internal/domain"
	"user-management/internal/models"
)

// swagger:parameters
// nolint:unused
type createParametes struct {
	// in: header
	headers struct {
		// user unique identifier
		UserOId string `json:"user_oid"`
	}
	// in: body
	body struct {
		// user password
		// required true
		Password string             `json:"password"`
		Article  models.ArticleData `json:"article"`
	}
}

func (a Article) Create(ctx controllers.Context) error {
	fmt.Println("> create article controller")
	formParams, err := ctx.FormParams()
	if err != nil {
		e := ctx.JSON(http.StatusBadRequest, controllers.ErrorBody{Error: "failed to parse form params"})
		return fmt.Errorf("%v: %w", e, err)
	}
	userOId := ctx.Request().Header.Get("user_oid")
	fmt.Println("> before exists")
	err = a.userModel.Exists(userOId, formParams.Get("password"))
	fmt.Print("> arter exists")
	if errors.Is(err, models.UserNotFoundErr) || errors.Is(err, models.InvalidAuthParameterErr) {
		e := ctx.JSON(http.StatusUnauthorized, controllers.ErrorBody{Error: err.Error()})
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
