package user_test

import (
	"fmt"
	"net/http"
	"testing"
	cntrl "user-management/internal/controllers"
	user "user-management/internal/controllers/login"
	cntrlMocks "user-management/internal/mocks/controllers"
	mdlMocks "user-management/internal/mocks/models"
	"user-management/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type loginFields struct {
	userModel models.UserModel
}
type loginArgs struct {
	ctx cntrl.Context
}
type loginMocks struct {
	userModel models.UserModel
	ctx       cntrl.Context
}

func TestLoginController_Register(t *testing.T) {
	type mockedRes struct {
		bindingErr    error
		createErr     error
		jsonCode      int
		noContentCode int
	}
	setup := func(mocks *loginMocks, res mockedRes) func() {
		ctx := mocks.ctx.(*cntrlMocks.Context).EXPECT()
		bindCall := ctx.Bind(mock.Anything).Return(res.bindingErr).Once()
		calls := []*mock.Call{
			bindCall,
			ctx.JSON(res.jsonCode, mock.Anything).Return(nil).NotBefore(bindCall).Maybe(),
			ctx.NoContent(res.noContentCode).Return(nil).NotBefore(bindCall).Maybe(),
			mocks.userModel.(*mdlMocks.UserModel).EXPECT().
				Create(mock.Anything).Return(res.createErr).NotBefore(bindCall).Once(),
		}
		return func() {
			for _, call := range calls {
				call.Unset()
			}
		}
	}
	ctx := cntrlMocks.NewContext(t)
	userModel := mdlMocks.NewUserModel(t)
	err := fmt.Errorf("invoked error")
	tests := []struct {
		name      string
		mockedRes mockedRes
		wantErr   error
	}{
		{
			name:      "binding failed",
			mockedRes: mockedRes{bindingErr: err, jsonCode: http.StatusBadRequest},
			wantErr:   err,
		},
		{
			name:      "username duplication",
			mockedRes: mockedRes{createErr: models.UsernameDuplicationErr, jsonCode: http.StatusConflict},
			wantErr:   models.UsernameDuplicationErr,
		},
		{
			name:      "invalid user credentials",
			mockedRes: mockedRes{createErr: models.InvalidAuthParameterErr, jsonCode: http.StatusBadRequest},
			wantErr:   models.InvalidAuthParameterErr,
		},
		{
			name:      "internal server error",
			mockedRes: mockedRes{createErr: err, noContentCode: http.StatusInternalServerError},
			wantErr:   err,
		},
		{
			name:      "success",
			mockedRes: mockedRes{noContentCode: http.StatusCreated},
			wantErr:   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanSetup := setup(&loginMocks{ctx: ctx, userModel: userModel}, tt.mockedRes)
			defer cleanSetup()
			err := user.NewLoginController(userModel).Register(ctx)
			if err != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			assert.Nil(t, err)
		})
	}
}

func TestLoginController_Authorize(t *testing.T) {
	type mockedRes struct {
		bindingErr    error
		authErr       error
		jsonCode      int
		noContentCode int
	}
	userId, token := "user-id", "token"
	setup := func(mocks *loginMocks, res mockedRes) func() {
		ctx := mocks.ctx.(*cntrlMocks.Context).EXPECT()
		bindCall := ctx.Bind(mock.Anything).Return(res.bindingErr).Once()
		calls := []*mock.Call{
			bindCall,
			ctx.JSON(res.jsonCode, mock.Anything).Return(nil).NotBefore(bindCall).Maybe(),
			ctx.NoContent(res.noContentCode).Return(nil).NotBefore(bindCall).Maybe(),
			mocks.userModel.(*mdlMocks.UserModel).EXPECT().
				Authorize(mock.Anything).Return(userId, token, res.authErr).NotBefore(bindCall).Once(),
		}
		return func() {
			for _, call := range calls {
				call.Unset()
			}
		}
	}
	ctx := cntrlMocks.NewContext(t)
	userModel := mdlMocks.NewUserModel(t)
	err := fmt.Errorf("invoked error")
	tests := []struct {
		name      string
		mockedRes mockedRes
		wantErr   error
	}{
		{
			name:      "binding failed",
			mockedRes: mockedRes{bindingErr: err, jsonCode: http.StatusUnauthorized},
			wantErr:   err,
		},
		{
			name:      "invalid credentials",
			mockedRes: mockedRes{authErr: models.InvalidAuthParameterErr, jsonCode: http.StatusUnauthorized},
			wantErr:   models.InvalidAuthParameterErr,
		},
		{
			name:      "internal server error",
			mockedRes: mockedRes{authErr: err, noContentCode: http.StatusInternalServerError},
			wantErr:   err,
		},
		{
			name:      "success",
			mockedRes: mockedRes{jsonCode: http.StatusOK},
			wantErr:   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanSetup := setup(&loginMocks{ctx: ctx, userModel: userModel}, tt.mockedRes)
			defer cleanSetup()
			err := user.NewLoginController(userModel).Authorize(ctx)
			if err != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			assert.Nil(t, err)
		})
	}
}
