package article

import (
	"fmt"
	"net/http"
	"user-management/internal/controllers"
	"user-management/internal/domain"
	"user-management/internal/models"
	"user-management/internal/tools"
)

// swagger:model
type articles struct {
	Data           []domain.Article      `json:"data"`
	PaginationData models.PaginationData `json:"pagination"`
}

// swagger:parameters articles_for_user
// nolint:unused
type articlesForUserParams struct {
	// username of owner of articles
	// in: query
	// type: string
	Username string `json:"username"`
	// number of page in pagination
	// in: query
	// type: "integer"
	// format: "int32"
	// minimum: 0
	// default: 0
	Page int `json:"page"`
	// maximal number of fetched articles
	// in: query
	// type: "integer"
	// format: "int32"
	// minimum: 0
	// maximum: 50
	// default: 10
	Limit int `json:"limit"`
}

// success
// swagger:response articlesForUserResp200
// nolint:unused
type articlesForUserResp200 struct {
	// in: body
	body articles
}

// invalid prameters provided
// swagger:response articlesForUserResp400
// nolint:unused
type articlesForUserResp400 struct {
	// in: body
	body controllers.ErrorBody
}

// swagger:route GET /articles articles articles_for_user
// get list of articles for specified user
// ---
// Validates `page` and `limit` params and returns list of articles for specified user by his `username`.
// New articles are in the start and old ones are in the end of list.
//
// responses:
//
//	200: articlesForUserResp200
//	400: articlesForUserResp400
//	500: commonError
func (a Article) GetForUser(ctx controllers.Context) error {
	query := ctx.QueryParams()
	username, pageStr, limitStr := query.Get("username"), query.Get("page"), query.Get("limit")
	page, limit, err := tools.NewPagination().Parse(pageStr, limitStr)
	if err != nil {
		e := ctx.JSON(http.StatusBadRequest, controllers.ErrorBody{Error: err.Error()})
		return fmt.Errorf("%w: %w", e, err)
	}

	articlesData, paginationData, err := a.articleModel.GetForUser(username, page, limit)
	if err != nil {
		e := ctx.NoContent(http.StatusInternalServerError)
		return fmt.Errorf("%w: %w", e, err)
	}

	return ctx.JSON(http.StatusOK, articles{Data: articlesData, PaginationData: paginationData})
}
