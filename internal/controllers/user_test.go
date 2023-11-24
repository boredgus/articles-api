package controllers

import (
	"fmt"
	"net/http"
	"testing"
	cntrlMocks "user-management/internal/mocks/controllers"
	mdlMocks "user-management/internal/mocks/models"
	"user-management/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// 	t.Run("authorization fails: user binding failed", func(t *testing.T) {
// 		bindMock := context.On("Bind", mock.Anything).Return(bindErr)
// 		jsonMock := context.On("JSON", http.StatusUnauthorized, mock.Anything).Return(nil)
// 		err := NewLoginController(model).Authorize(context)
// 		assert.ErrorIs(t, err, bindErr)
// 		bindMock.Unset()
// 		jsonMock.Unset()
// 	})

// 	t.Run("authorization fails if user credentials are invalid", func(t *testing.T) {
// 		bindMock := context.On("Bind", mock.Anything).Return(nil)
// 		modelMock := model.On("Authorize", mock.Anything).Return("", "", models.InvalidAuthParameterErr)
// 		jsonMock := context.On("JSON", http.StatusUnauthorized, mock.Anything).Return(nil)
// 		err := NewLoginController(model).Authorize(context)
// 		assert.ErrorIs(t, err, models.InvalidAuthParameterErr)
// 		bindMock.Unset()
// 		modelMock.Unset()
// 		jsonMock.Unset()
// 	})

// 	t.Run("authorization fails because of internal error", func(t *testing.T) {
// 		bindMock := context.On("Bind", mock.Anything).Return(nil)
// 		modelMock := model.On("Authorize", mock.Anything).Return("", "", serverError)
// 		jsonMock := context.On("NoContent", http.StatusInternalServerError, mock.Anything).Return(nil)
// 		err := NewLoginController(model).Authorize(context)
// 		assert.ErrorIs(t, err, serverError)
// 		bindMock.Unset()
// 		modelMock.Unset()
// 		jsonMock.Unset()
// 	})

// 	t.Run("successful authorization", func(t *testing.T) {
// 		bindMock := context.On("Bind", mock.Anything).Return(nil)
// 		userId, token := "user-id", "token"
// 		modelMock := model.On("Authorize", mock.Anything).Return(userId, token, nil)
// 		jsonMock := context.On("JSON", http.StatusOK, AuthBody{Token: token, UserId: userId}).Return(nil)
// 		err := NewLoginController(model).Authorize(context)
// 		assert.Nil(t, err)
// 		bindMock.Unset()
// 		modelMock.Unset()
// 		jsonMock.Unset()
// 	})
// }

type loginFields struct {
	userModel models.UserModel
}
type loginArgs struct {
	ctx Context
}
type loginMocks struct {
	userModel models.UserModel
	ctx       Context
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
			err := NewLoginController(tt.fields.userModel).Register(tt.args.ctx)
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
			err := NewLoginController(tt.fields.userModel).Authorize(tt.args.ctx)
			if err != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			assert.Nil(t, err)
		})
	}
}
