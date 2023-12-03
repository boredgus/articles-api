package article

import (
	"errors"
	"net/http"
	"net/url"
	"testing"
	cntrs "user-management/internal/controllers"
	"user-management/internal/domain"
	cntlrMocks "user-management/internal/mocks/controllers"
	mdlMocks "user-management/internal/mocks/models"
	mdl "user-management/internal/models"
	"user-management/internal/views"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestArticleController_Create(t *testing.T) {
	type mockedRes struct {
		jsonCode      int
		jsonBody      interface{}
		noContentCode int
		userExistErr  error
		bindErr       error
		createArticle domain.Article
		createErr     error
	}
	ctxMock := cntlrMocks.NewContext(t)
	articleModelMock := mdlMocks.NewArticleModel(t)
	userModelMock := mdlMocks.NewUserModel(t)
	setup := func(res mockedRes) func() {
		h := http.Header{}
		h.Set(UserOIdKey, "user-identifier")
		requestCall := ctxMock.EXPECT().
			Request().Return(&http.Request{Header: h}).Once()
		q := url.Values{}
		q.Set("passwrod", "pass")
		queryParamsCall := ctxMock.EXPECT().QueryParams().NotBefore(requestCall).Return(q).Once()
		existsCall := userModelMock.EXPECT().
			Exists(h.Get(UserOIdKey), q.Get("password")).NotBefore(requestCall, queryParamsCall).
			Return(res.userExistErr).Once()
		calls := []*mock.Call{
			requestCall, queryParamsCall, existsCall,
			ctxMock.EXPECT().
				JSON(res.jsonCode, res.jsonBody).Return(nil).
				NotBefore(requestCall, queryParamsCall, existsCall).Maybe(),
			ctxMock.EXPECT().
				NoContent(res.noContentCode).Return(nil).
				NotBefore(requestCall, queryParamsCall, existsCall).Maybe(),
			ctxMock.EXPECT().Bind(mock.Anything).
				NotBefore(requestCall, queryParamsCall, existsCall).Return(res.bindErr).Maybe(),
			articleModelMock.EXPECT().
				Create(h.Get(UserOIdKey), &res.createArticle).NotBefore(requestCall, existsCall).
				Return(res.createErr).Maybe(),
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
			name: "invalid user_oid or password",
			mockedRes: mockedRes{
				userExistErr: mdl.InvalidAuthParameterErr,
				jsonCode:     http.StatusUnauthorized,
				jsonBody:     cntrs.ErrorBody{Error: "invalid user_oid or password"},
			},
			wantErr: mdl.InvalidAuthParameterErr,
		},
		{
			name: "failed to get user data",
			mockedRes: mockedRes{
				userExistErr:  someError,
				noContentCode: http.StatusInternalServerError,
			},
			wantErr: someError,
		},
		{
			name: "failed to get article data",
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
				createErr:     mdl.InvalidArticleErr,
				jsonCode:      http.StatusBadRequest,
				jsonBody:      cntrs.ErrorBody{Error: mdl.InvalidArticleErr.Error()},
			},
			wantErr: mdl.InvalidArticleErr,
		},
		{
			name: "failed to create article",
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
			err := NewArticleController(userModelMock, articleModelMock).Create(ctxMock)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			assert.Nil(t, err)
		})
	}
}
