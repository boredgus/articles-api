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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestArticleController_Create(t *testing.T) {
	type mockedRes struct {
		jsonCode      int
		jsonBody      interface{}
		noContentCode int
		formParams    url.Values
		formParamsErr error
		userExistErr  error
		createUserOId string
		createArticle domain.Article
		createErr     error
	}
	ctxMock := cntlrMocks.NewContext(t)
	articleModelMock := mdlMocks.NewArticleModel(t)
	userModelMock := mdlMocks.NewUserModel(t)
	setup := func(res mockedRes) func() {
		formParamsCall := ctxMock.EXPECT().
			FormParams().Return(res.formParams, res.formParamsErr).Once()
		requestCall := ctxMock.EXPECT().
			Request().Return(&http.Request{Header: http.Header{}}).
			NotBefore(formParamsCall).Maybe()
		existsCall := userModelMock.EXPECT().
			Exists(mock.Anything, mock.Anything).NotBefore(formParamsCall, requestCall).
			Return(res.userExistErr).Maybe()
		calls := []*mock.Call{
			formParamsCall,
			requestCall,
			ctxMock.EXPECT().
				JSON(res.jsonCode, res.jsonBody).Return(nil).NotBefore(formParamsCall).Maybe(),
			existsCall,
			ctxMock.EXPECT().
				NoContent(res.noContentCode).Return(nil).NotBefore(formParamsCall).Maybe(),
			articleModelMock.EXPECT().
				Create(res.createUserOId, &res.createArticle).NotBefore(formParamsCall, requestCall, existsCall).
				Return(res.createErr).Maybe(),
		}
		return func() {
			for _, call := range calls {
				call.Unset()
			}
		}
	}
	someError := errors.New("some error")
	tests := []struct {
		name      string
		mockedRes mockedRes
		wantErr   error
	}{
		{
			name: "failed to get form params",
			mockedRes: mockedRes{
				formParamsErr: someError,
				jsonCode:      http.StatusBadRequest,
				jsonBody:      mock.Anything,
			},
			wantErr: someError,
		},
		{
			name: "invalid user oid",
			mockedRes: mockedRes{
				userExistErr: mdl.UserNotFoundErr,
				jsonCode:     http.StatusUnauthorized,
				jsonBody:     cntrs.ErrorBody{Error: mdl.UserNotFoundErr.Error()},
			},
			wantErr: mdl.UserNotFoundErr,
		},
		{
			name: "invalid user password",
			mockedRes: mockedRes{
				userExistErr: mdl.InvalidAuthParameterErr,
				jsonCode:     http.StatusUnauthorized,
				jsonBody:     cntrs.ErrorBody{Error: mdl.InvalidAuthParameterErr.Error()},
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
			name: "invalid article data provided",
			mockedRes: mockedRes{
				createErr: mdl.InvalidArticleErr,
				jsonCode:  http.StatusBadRequest,
				jsonBody:  cntrs.ErrorBody{Error: mdl.InvalidArticleErr.Error()},
			},
			wantErr: mdl.InvalidArticleErr,
		},
		{
			name: "failed to create article",
			mockedRes: mockedRes{
				createErr:     someError,
				noContentCode: http.StatusInternalServerError,
			},
			wantErr: someError,
		},
		{
			name: "success",
			mockedRes: mockedRes{
				jsonCode: http.StatusCreated,
				jsonBody: domain.Article{},
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
