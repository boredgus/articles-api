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
	valueOfFields := loginFields{userModel: mdlMocks.NewUserModel(t)}
	valueOfArgs := loginArgs{ctx: cntrlMocks.NewContext(t)}
	err := fmt.Errorf("invoked error")
	tests := []struct {
		name      string
		fields    loginFields
		mockedRes mockedRes
		args      loginArgs
		wantErr   error
	}{
		{
			name:      "binding failed",
			fields:    valueOfFields,
			mockedRes: mockedRes{bindingErr: err, jsonCode: http.StatusBadRequest},
			args:      valueOfArgs,
			wantErr:   err,
		},
		{
			name:      "username duplication",
			fields:    valueOfFields,
			mockedRes: mockedRes{createErr: models.UsernameDuplicationErr, jsonCode: http.StatusConflict},
			args:      valueOfArgs,
			wantErr:   models.UsernameDuplicationErr,
		},
		{
			name:      "invalid user credentials",
			fields:    valueOfFields,
			mockedRes: mockedRes{createErr: models.InvalidAuthParameterErr, jsonCode: http.StatusBadRequest},
			args:      valueOfArgs,
			wantErr:   models.InvalidAuthParameterErr,
		},
		{
			name:      "internal server error",
			fields:    valueOfFields,
			mockedRes: mockedRes{createErr: err, noContentCode: http.StatusInternalServerError},
			args:      valueOfArgs,
			wantErr:   err,
		},
		{
			name:      "success",
			fields:    valueOfFields,
			mockedRes: mockedRes{noContentCode: http.StatusCreated},
			args:      valueOfArgs,
			wantErr:   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanSetup := setup(&loginMocks{ctx: tt.args.ctx, userModel: tt.fields.userModel}, tt.mockedRes)
			defer cleanSetup()
			err := user.NewLoginController(tt.fields.userModel).Register(tt.args.ctx)
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
	valueOfFields := loginFields{userModel: mdlMocks.NewUserModel(t)}
	valueOfArgs := loginArgs{ctx: cntrlMocks.NewContext(t)}
	err := fmt.Errorf("invoked error")
	tests := []struct {
		name      string
		fields    loginFields
		mockedRes mockedRes
		args      loginArgs
		wantErr   error
	}{
		{
			name:      "binding failed",
			fields:    valueOfFields,
			mockedRes: mockedRes{bindingErr: err, jsonCode: http.StatusUnauthorized},
			args:      valueOfArgs,
			wantErr:   err,
		},
		{
			name:      "invalid credentials",
			fields:    valueOfFields,
			mockedRes: mockedRes{authErr: models.InvalidAuthParameterErr, jsonCode: http.StatusUnauthorized},
			args:      valueOfArgs,
			wantErr:   models.InvalidAuthParameterErr,
		},
		{
			name:      "internal server error",
			fields:    valueOfFields,
			mockedRes: mockedRes{authErr: err, noContentCode: http.StatusInternalServerError},
			args:      valueOfArgs,
			wantErr:   err,
		},
		{
			name:      "success",
			fields:    valueOfFields,
			mockedRes: mockedRes{jsonCode: http.StatusOK},
			args:      valueOfArgs,
			wantErr:   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanSetup := setup(&loginMocks{ctx: tt.args.ctx, userModel: tt.fields.userModel}, tt.mockedRes)
			defer cleanSetup()
			err := user.NewLoginController(tt.fields.userModel).Authorize(tt.args.ctx)
			if err != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			assert.Nil(t, err)
		})
	}
}
