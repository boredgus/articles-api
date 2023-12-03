package article

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"user-management/internal/controllers"
	"user-management/internal/domain"
	cntlrMocks "user-management/internal/mocks/controllers"
	mdlMocks "user-management/internal/mocks/models"
	"user-management/internal/models"
	"user-management/internal/tools"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestArticle_GetForUser(t *testing.T) {
	type params struct {
		username string
		page     int
		limit    int
	}
	type fetchRes struct {
		articles   []domain.Article
		pagination models.PaginationData
		err        error
	}
	type mockedRes struct {
		jsonCode      int
		jsonBody      interface{}
		noContentCode int
		queryParams   url.Values
		fetchParams   params
		fetchRes      fetchRes
	}
	ctxMock := cntlrMocks.NewContext(t)
	articleModelMock := mdlMocks.NewArticleModel(t)
	setup := func(res mockedRes) func() {
		queryParamsCall := ctxMock.EXPECT().QueryParams().Return(res.queryParams).Once()
		calls := []*mock.Call{
			queryParamsCall,
			ctxMock.EXPECT().
				JSON(res.jsonCode, mock.Anything).Return(nil).NotBefore(queryParamsCall).Maybe(),
			ctxMock.EXPECT().
				NoContent(res.noContentCode).Return(nil).NotBefore(queryParamsCall).Maybe(),
			articleModelMock.EXPECT().
				GetForUser(res.fetchParams.username, res.fetchParams.page, res.fetchParams.limit).NotBefore(queryParamsCall).
				Return(res.fetchRes.articles, res.fetchRes.pagination, res.fetchRes.err).Once(),
		}
		return func() {
			for _, call := range calls {
				call.Unset()
			}
		}
	}
	valildParams := params{username: "user", page: 1, limit: 8}
	validQueryParams := url.Values{
		"username": []string{valildParams.username},
		"page":     []string{fmt.Sprint(valildParams.page)},
		"limit":    []string{fmt.Sprint(valildParams.limit)},
	}
	articls := []domain.Article{
		{Theme: "first", Tags: []string{}},
		{Theme: "second", Tags: []string{"update"}},
	}
	someErr := errors.New("some err")
	tests := []struct {
		name      string
		mockedRes mockedRes
		wantErr   error
	}{
		{
			name: "pagination data is invalid",
			mockedRes: mockedRes{
				queryParams: url.Values{"page": []string{"-1"}},
				jsonCode:    http.StatusBadRequest,
				jsonBody:    controllers.ErrorBody{Error: tools.PageOutOfRangeErr.Error()},
			},
			wantErr: tools.PageOutOfRangeErr,
		},
		{
			name: "failed to fetch articles",
			mockedRes: mockedRes{
				queryParams:   validQueryParams,
				fetchParams:   valildParams,
				fetchRes:      fetchRes{err: someErr},
				noContentCode: http.StatusInternalServerError,
			},
			wantErr: someErr,
		},
		{
			name: "success",
			mockedRes: mockedRes{
				queryParams: validQueryParams,
				fetchParams: valildParams,
				fetchRes:    fetchRes{articles: articls, pagination: models.PaginationData{Count: len(articls)}},
				jsonCode:    http.StatusOK,
				jsonBody:    articles{PaginationData: models.PaginationData{Limit: valildParams.limit}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanSetup := setup(tt.mockedRes)
			defer cleanSetup()
			err := NewArticleController(mdlMocks.NewUserModel(t), articleModelMock).
				GetForUser(ctxMock)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			assert.Nil(t, err)
		})
	}
}
