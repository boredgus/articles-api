package article

import (
	"a-article/internal/auth"
	cntrs "a-article/internal/controllers"
	"a-article/internal/domain"
	cntlrMocks "a-article/internal/mocks/controllers"
	mdlMocks "a-article/internal/mocks/models"
	mdl "a-article/internal/models"
	"a-article/internal/views"
	"errors"
	"net/http"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestArticleController_Create(t *testing.T) {
	type mockedRes struct {
		jsonCode      int
		jsonBody      interface{}
		noContentCode int
		bindErr       error
		createArticle domain.Article
		createErr     error
	}
	ctxMock := cntlrMocks.NewContext(t)
	articleModelMock := mdlMocks.NewArticleModel(t)
	setup := func(res mockedRes) func() {
		bindCall := ctxMock.EXPECT().Bind(mock.Anything).Return(res.bindErr).Once()
		userOId := "user-oid"
		getCall := ctxMock.EXPECT().Get("user").NotBefore(bindCall).
			Return(jwt.NewWithClaims(jwt.SigningMethodHS256,
				&auth.JWTClaims{JWTPayload: auth.JWTPayload{UserOId: userOId}})).Once()
		calls := []*mock.Call{
			bindCall, getCall,
			articleModelMock.EXPECT().
				Create(userOId, &res.createArticle).NotBefore().
				Return(res.createErr).Maybe(),
			ctxMock.EXPECT().
				JSON(res.jsonCode, res.jsonBody).Return(nil).
				NotBefore(bindCall).Maybe(),
			ctxMock.EXPECT().
				NoContent(res.noContentCode).Return(nil).
				NotBefore().Maybe(),
		}
		return func() {
			for _, call := range calls {
				call.Unset()
			}
		}
	}
	someError := errors.New("some error")
	artcl := domain.Article{Tags: []string{}}
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
				jsonBody: cntrs.ErrorBody{Error: "failed to parse article"},
			},
			wantErr: someError,
		},
		{
			name: "invalid article data provided",
			mockedRes: mockedRes{
				createArticle: artcl,
				createErr:     mdl.InvalidDataErr,
				jsonCode:      http.StatusBadRequest,
				jsonBody:      cntrs.ErrorBody{Error: mdl.InvalidDataErr.Error()},
			},
			wantErr: mdl.InvalidDataErr,
		},
		{
			name: "internal server error",
			mockedRes: mockedRes{
				createArticle: artcl,
				createErr:     someError,
				noContentCode: http.StatusInternalServerError,
			},
			wantErr: someError,
		},
		{
			name: "success",
			mockedRes: mockedRes{
				createArticle: artcl,
				jsonCode:      http.StatusCreated,
				jsonBody:      views.NewArticleView(artcl),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanSetup := setup(tt.mockedRes)
			defer cleanSetup()
			err := NewArticleController(articleModelMock).Create(ctxMock)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			assert.Nil(t, err)
		})
	}
}
