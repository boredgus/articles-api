package user

import (
	"a-article/internal/auth"
	"a-article/internal/controllers"
	"a-article/internal/models"
	"fmt"
	"net/http"
	"testing"

	cntrlMocks "a-article/internal/mocks/controllers"
	mdlMocks "a-article/internal/mocks/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserController_UpdateRole(t *testing.T) {
	type mockedRes struct {
		bindErr       error
		updateErr     error
		jsonCode      int
		jsonBody      any
		noContentCode int
	}
	ctxMock := cntrlMocks.NewContext(t)
	userModelMock := mdlMocks.NewUserModel(t)
	issuerRole, userId := "admin", "oid"
	setup := func(res mockedRes) func() {
		bindCall := ctxMock.EXPECT().Bind(mock.Anything).Return(res.bindErr).Once()
		getCall := ctxMock.EXPECT().Get("user").
			Return(jwt.NewWithClaims(jwt.SigningMethodHS256,
				&auth.JWTClaims{JWTPayload: auth.JWTPayload{Role: issuerRole}})).Once()
		pathParamCall := ctxMock.EXPECT().PathParam("user_id").NotBefore(getCall).Return(userId)
		updateCall := userModelMock.EXPECT().
			UpdateRole(issuerRole, userId, "").Return(res.updateErr).NotBefore(pathParamCall).Once()
		calls := []*mock.Call{
			bindCall, getCall, pathParamCall, updateCall,
			ctxMock.EXPECT().JSON(res.jsonCode, res.jsonBody).Return(nil).NotBefore(bindCall).Maybe(),
			ctxMock.EXPECT().NoContent(res.noContentCode).Return(nil).NotBefore(updateCall).Maybe(),
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
			name: "failed to bind role",
			mockedRes: mockedRes{
				bindErr:  someErr,
				jsonCode: http.StatusBadRequest,
				jsonBody: controllers.ErrorBody{Error: "failed to bind data"},
			},
			wantErr: someErr,
		},
		{
			name: "invalid role value",
			mockedRes: mockedRes{
				updateErr: models.InvalidDataErr,
				jsonCode:  http.StatusBadRequest,
				jsonBody:  controllers.ErrorBody{Error: models.InvalidDataErr.Error()},
			},
			wantErr: models.InvalidDataErr,
		},
		{
			name: "user not found",
			mockedRes: mockedRes{
				updateErr: models.UserNotFoundErr,
				jsonCode:  http.StatusNotFound,
				jsonBody:  controllers.ErrorBody{Error: models.UserNotFoundErr.Error()},
			},
			wantErr: models.UserNotFoundErr,
		},
		{
			name: "not enough rights to update user",
			mockedRes: mockedRes{
				updateErr: models.NotEnoughRightsErr,
				jsonCode:  http.StatusForbidden,
				jsonBody:  controllers.ErrorBody{Error: models.NotEnoughRightsErr.Error()},
			},
			wantErr: models.NotEnoughRightsErr,
		},
		{
			name: "internal server error",
			mockedRes: mockedRes{
				updateErr:     someErr,
				noContentCode: http.StatusInternalServerError,
			},
			wantErr: someErr,
		},
		{
			name: "success",
			mockedRes: mockedRes{
				jsonCode: http.StatusOK,
				jsonBody: controllers.InfoResponse{Message: "successfuly updated role to "},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanSetup := setup(tt.mockedRes)
			defer cleanSetup()
			err := NewUserController(userModelMock).UpdateRole(ctxMock)
			if err != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			assert.Nil(t, err)
		})
	}
}
