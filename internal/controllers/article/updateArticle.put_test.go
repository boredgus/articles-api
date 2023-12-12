package article

import (
	"errors"
	"net/http"
	"testing"
	"user-management/internal/auth"
	"user-management/internal/controllers"
	"user-management/internal/domain"
	cntlrMocks "user-management/internal/mocks/controllers"
	mdlMocks "user-management/internal/mocks/models"
	"user-management/internal/models"
	"user-management/internal/views"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestArticleController_Update(t *testing.T) {
	type mockedRes struct {
		bindErr       error
		userRole      string
		jsonCode      int
		jsonBody      interface{}
		noContentCode int
		updateArticle domain.Article
		updateErr     error
	}
	ctxMock := cntlrMocks.NewContext(t)
	articleModelMock := mdlMocks.NewArticleModel(t)
	userModelMock := mdlMocks.NewUserModel(t)
	articleId := "artice-id"
	setup := func(res mockedRes) func() {
		bindCall := ctxMock.EXPECT().Bind(mock.Anything).Return(res.bindErr).Maybe()
		userOId := "user-oid"
		pathParamCall := ctxMock.EXPECT().
			PathParam("article_id").Return(articleId).Once().NotBefore(bindCall)
		getCall := ctxMock.EXPECT().Get("user").NotBefore(bindCall).
			Return(jwt.NewWithClaims(jwt.SigningMethodHS256,
				&auth.JWTClaims{JWTPayload: auth.JWTPayload{UserOId: userOId, Role: res.userRole}})).Once()
		updateCall := articleModelMock.EXPECT().
			Update(userOId, res.userRole, &res.updateArticle).NotBefore(getCall).
			Return(res.updateErr).Maybe()
		calls := []*mock.Call{
			bindCall,
			pathParamCall,
			getCall,
			updateCall,
			ctxMock.EXPECT().
				JSON(res.jsonCode, res.jsonBody).Return(nil).
				NotBefore(bindCall).Maybe(),
			ctxMock.EXPECT().
				NoContent(res.noContentCode).Return(nil).
				NotBefore(updateCall).Maybe(),
		}
		return func() {
			for _, call := range calls {
				call.Unset()
			}
		}
	}
	someError := errors.New("some error")
	artcl := domain.Article{OId: articleId, Tags: []string{}}
	tests := []struct {
		name      string
		mockedRes mockedRes
		wantErr   error
	}{
		{
			name: "failed to bind article data",
			mockedRes: mockedRes{
				bindErr:  someError,
				jsonCode: http.StatusBadRequest,
				jsonBody: controllers.ErrorBody{Error: "failed to parse article"}},
			wantErr: someError,
		},
		{
			name: "article with such oid does not exist",
			mockedRes: mockedRes{
				updateArticle: artcl,
				updateErr:     models.ArticleNotFoundErr,
				jsonCode:      http.StatusNotFound,
				jsonBody:      mock.Anything},
			wantErr: models.ArticleNotFoundErr,
		},
		{
			name: "user does not have rights to change this article",
			mockedRes: mockedRes{
				updateArticle: artcl,
				updateErr:     models.NotEnoughRightsErr,
				jsonCode:      http.StatusForbidden,
				jsonBody:      mock.Anything},
			wantErr: models.NotEnoughRightsErr,
		},
		{
			name: "article data is invalid",
			mockedRes: mockedRes{
				updateArticle: artcl,
				updateErr:     models.InvalidArticleErr,
				jsonCode:      http.StatusBadRequest,
				jsonBody:      mock.Anything},
			wantErr: models.InvalidArticleErr,
		},
		{
			name: "internal server error",
			mockedRes: mockedRes{
				updateArticle: artcl,
				updateErr:     someError,
				noContentCode: http.StatusInternalServerError},
			wantErr: someError,
		},
		{
			name: "success",
			mockedRes: mockedRes{
				updateArticle: artcl,
				jsonCode:      http.StatusOK,
				jsonBody:      views.NewArticleView(artcl)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanSetup := setup(tt.mockedRes)
			defer cleanSetup()
			err := NewArticleController(userModelMock, articleModelMock).Update(ctxMock)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			assert.Nil(t, err)
		})
	}
}
