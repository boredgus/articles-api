package article

import (
	"errors"
	"net/http"
	"testing"
	"user-management/internal/auth"
	cntlrMocks "user-management/internal/mocks/controllers"
	mdlMocks "user-management/internal/mocks/models"
	"user-management/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestArticle_Delete(t *testing.T) {
	type mockedRes struct {
		userRole      string
		jsonCode      int
		jsonBody      interface{}
		noContentCode int
		deleteErr     error
	}
	ctxMock := cntlrMocks.NewContext(t)
	articleModelMock := mdlMocks.NewArticleModel(t)
	articleId := "artice-id"
	setup := func(res mockedRes) func() {
		userOId := "user-oid"
		getCall := ctxMock.EXPECT().Get("user").
			Return(jwt.NewWithClaims(jwt.SigningMethodHS256,
				&auth.JWTClaims{JWTPayload: auth.JWTPayload{UserOId: userOId, Role: res.userRole}})).Once()
		pathParamCall := ctxMock.EXPECT().
			PathParam("article_id").Return(articleId).Once().NotBefore(getCall)
		deleteCall := articleModelMock.EXPECT().
			Delete(userOId, res.userRole, articleId).NotBefore(pathParamCall).
			Return(res.deleteErr).Maybe()
		calls := []*mock.Call{
			getCall,
			pathParamCall,
			deleteCall,
			ctxMock.EXPECT().
				JSON(res.jsonCode, res.jsonBody).Return(nil).
				NotBefore(deleteCall).Maybe(),
			ctxMock.EXPECT().
				NoContent(res.noContentCode).Return(nil).
				NotBefore(deleteCall).Maybe(),
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
			name: "article with such oid does not exists",
			mockedRes: mockedRes{
				deleteErr: models.NotFoundErr,
				jsonCode:  http.StatusNotFound,
				jsonBody:  mock.Anything,
			},
			wantErr: models.NotFoundErr,
		},
		{
			name: "not enough rights to delete article",
			mockedRes: mockedRes{
				deleteErr: models.NotEnoughRightsErr,
				jsonCode:  http.StatusForbidden,
				jsonBody:  mock.Anything,
			},
			wantErr: models.NotEnoughRightsErr,
		},
		{
			name: "internal server error",
			mockedRes: mockedRes{
				deleteErr:     someError,
				noContentCode: http.StatusInternalServerError,
			},
			wantErr: someError,
		},
		{
			name: "success",
			mockedRes: mockedRes{
				noContentCode: http.StatusOK,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanSetup := setup(tt.mockedRes)
			defer cleanSetup()
			err := NewArticleController(articleModelMock).Delete(ctxMock)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			assert.Nil(t, err)
		})
	}
}
