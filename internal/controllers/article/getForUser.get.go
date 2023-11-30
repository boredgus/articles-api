package article

import (
	"fmt"
	"net/http"
	"strconv"
	"user-management/internal/controllers"
	"user-management/internal/domain"
	"user-management/internal/models"
)

// swagger: model
type articles struct {
	Data           []domain.Article      `json:"data"`
	PaginationData models.PaginationData `json:"pagination"`
}

// swagger:response articlesForUserResp200
// nolint:unused
type articlesForUserResp200 struct {
	// in: body
	body articles
}

const DefaultPage = 0
const DefaultLimit = 10

func (a Article) GetForUser(ctx controllers.Context) error {
	username, pageStr, limitStr := ctx.QueryParams().Get("username"), ctx.QueryParams().Get("page"), ctx.QueryParams().Get("limit")
	if len(pageStr) == 0 {
		pageStr = "0"
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		e := ctx.JSON(http.StatusBadRequest, controllers.ErrorBody{Error: "page should be number"})
		return fmt.Errorf("%w: %w", e, err)
	}
	if len(limitStr) == 0 {
		limitStr = "10"
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		e := ctx.JSON(http.StatusBadRequest, controllers.ErrorBody{Error: "limit should be number"})
		return fmt.Errorf("%w: %w", e, err)
	}

	articlesData, paginationData, err := a.articleModel.GetForUser(username, page, limit)
	if err != nil {
		e := ctx.JSON(http.StatusBadRequest, controllers.ErrorBody{Error: "failed to get list of articles for user"})
		return fmt.Errorf("%w: %w", e, err)
	}

	return ctx.JSON(http.StatusOK, articles{Data: articlesData, PaginationData: paginationData})
}
