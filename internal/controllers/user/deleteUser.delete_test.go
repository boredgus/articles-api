package user_test

import (
	"a-article/internal/auth"
	user "a-article/internal/controllers/user"
	cntrlMocks "a-article/internal/mocks/controllers"
	mdlMocks "a-article/internal/mocks/models"
	"a-article/internal/models"
	"fmt"
	"net/http"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserController_Delete(t *testing.T) {
	type mockedRes struct {
		deleteErr     error
		jsonCode      int
		noContentCode int
	}
	ctxMock := cntrlMocks.NewContext(t)
	userModelMock := mdlMocks.NewUserModel(t)
	setup := func(res mockedRes) func() {
		userId, role := "oid", "role"
		getCall := ctxMock.EXPECT().Get("user").
			Return(jwt.NewWithClaims(jwt.SigningMethodHS256,
				&auth.JWTClaims{JWTPayload: auth.JWTPayload{Role: role}})).Once()
		pathParamCall := ctxMock.EXPECT().PathParam("user_id").NotBefore(getCall).Return(userId)
		deleteCall := userModelMock.EXPECT().
			Delete(mock.Anything, mock.Anything).Return(res.deleteErr).NotBefore(pathParamCall).Once()
		calls := []*mock.Call{
			getCall, pathParamCall, deleteCall,
			ctxMock.EXPECT().JSON(res.jsonCode, mock.Anything).Return(nil).NotBefore(deleteCall).Maybe(),
			ctxMock.EXPECT().NoContent(res.noContentCode).Return(nil).NotBefore(deleteCall).Maybe(),
		}
		return func() {
			for _, call := range calls {
				call.Unset()
			}
		}
	}
	someErr := fmt.Errorf("invoked error")
	tests := []struct {
		name      string
		mockedRes mockedRes
		wantErr   error
	}{
		{
			name: "user with provided oid does not exists",
			mockedRes: mockedRes{
				deleteErr: models.UserNotFoundErr,
				jsonCode:  http.StatusNotFound,
			},
			wantErr: models.UserNotFoundErr,
		},
		{
			name: "not enough rights to delete user",
			mockedRes: mockedRes{
				deleteErr: models.NotEnoughRightsErr,
				jsonCode:  http.StatusForbidden,
			},
			wantErr: models.NotEnoughRightsErr,
		},
		{
			name: "not enough rights to delete user",
			mockedRes: mockedRes{
				deleteErr:     someErr,
				noContentCode: http.StatusInternalServerError,
			},
			wantErr: someErr,
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
			err := user.NewUserController(userModelMock).Delete(ctxMock)
			if err != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			assert.Nil(t, err)
		})
	}
}
