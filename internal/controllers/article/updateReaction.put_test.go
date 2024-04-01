package article

import (
	"a-article/internal/auth"
	cntlrMocks "a-article/internal/mocks/controllers"
	mdlMocks "a-article/internal/mocks/models"
	"a-article/internal/models"
	"errors"
	"net/http"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestArticleController_UpdateReactionForArticle(t *testing.T) {
	type mockedRes struct {
		bindErr       error
		jsonCode      int
		noContentCode int
		updateErr     error
	}
	ctxMock := cntlrMocks.NewContext(t)
	articleModelMock := mdlMocks.NewArticleModel(t)
	setup := func(res mockedRes) func() {
		bindCall := ctxMock.EXPECT().Bind(mock.Anything).Return(res.bindErr).Maybe()
		userOId, articleId := "user-oid", "artice-id"
		getCall := ctxMock.EXPECT().Get("user").NotBefore(bindCall).
			Return(jwt.NewWithClaims(jwt.SigningMethodHS256,
				&auth.JWTClaims{JWTPayload: auth.JWTPayload{UserOId: userOId}})).Once()
		pathParamCall := ctxMock.EXPECT().
			PathParam("article_id").Return(articleId).Once().NotBefore(bindCall)
		updateCall := articleModelMock.EXPECT().
			UpdateReaction(userOId, articleId, "").NotBefore(pathParamCall).
			Return(res.updateErr).Maybe()
		calls := []*mock.Call{
			bindCall, getCall, pathParamCall, updateCall,
			ctxMock.EXPECT().
				JSON(res.jsonCode, mock.Anything).Return(nil).NotBefore(bindCall).Maybe(),
			ctxMock.EXPECT().
				NoContent(res.noContentCode).Return(nil).NotBefore(updateCall).Maybe(),
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
			name: "failed to bind payload",
			mockedRes: mockedRes{
				bindErr:  someError,
				jsonCode: http.StatusBadRequest,
			},
			wantErr: someError,
		},
		{
			name: "article with such id does not exists",
			mockedRes: mockedRes{
				updateErr: models.NotFoundErr,
				jsonCode:  http.StatusNotFound,
			},
			wantErr: models.NotFoundErr,
		},
		{
			name: "user is forbidden to update reaction",
			mockedRes: mockedRes{
				updateErr: models.NotEnoughRightsErr,
				jsonCode:  http.StatusForbidden,
			},
			wantErr: models.NotEnoughRightsErr,
		},
		{
			name: "invalid reaction provided",
			mockedRes: mockedRes{
				updateErr: models.InvalidDataErr,
				jsonCode:  http.StatusBadRequest,
			},
			wantErr: models.InvalidDataErr,
		},
		{
			name: "internal server error",
			mockedRes: mockedRes{
				updateErr:     someError,
				noContentCode: http.StatusInternalServerError},
			wantErr: someError,
		},
		{
			name:      "success",
			mockedRes: mockedRes{jsonCode: http.StatusOK},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanSetup := setup(tt.mockedRes)
			defer cleanSetup()
			err := NewArticleController(articleModelMock).UpdateReactionForArticle(ctxMock)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			assert.Nil(t, err)
		})
	}
}
