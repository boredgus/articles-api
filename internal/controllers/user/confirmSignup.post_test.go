package user_test

import (
	user "a-article/internal/controllers/user"
	cntrlMocks "a-article/internal/mocks/controllers"
	mdlMocks "a-article/internal/mocks/models"
	"a-article/internal/models"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserController_ConfirmSignup(t *testing.T) {
	type mockedRes struct {
		modelCallErr  error
		jsonCode      int
		noContentCode int
	}
	ctxMock := cntrlMocks.NewContext(t)
	userModelMock := mdlMocks.NewUserModel(t)
	setup := func(res mockedRes) func() {
		formCall := ctxMock.EXPECT().FormValue(mock.Anything).Return("").Twice()
		modelCall := userModelMock.EXPECT().ConfirmSignup("", "").Return(res.modelCallErr).NotBefore(formCall).Once()
		calls := []*mock.Call{
			formCall,
			modelCall,
			ctxMock.EXPECT().JSON(res.jsonCode, mock.Anything).Return(nil).NotBefore(modelCall).Maybe(),
			ctxMock.EXPECT().NoContent(res.noContentCode).Return(nil).NotBefore(formCall).Maybe(),
		}
		return func() {
			for _, call := range calls {
				call.Unset()
			}
		}
	}
	err := fmt.Errorf("invoked error")
	tests := []struct {
		name      string
		mockedRes mockedRes
		wantErr   error
	}{
		{
			name: "username duplicatiton",
			mockedRes: mockedRes{
				modelCallErr: models.UsernameDuplicationErr,
				jsonCode:     http.StatusConflict},
			wantErr: models.UsernameDuplicationErr,
		},
		{
			name: "related signup request not found",
			mockedRes: mockedRes{
				modelCallErr: models.NotFoundErr,
				jsonCode:     http.StatusNotFound},
			wantErr: models.NotFoundErr,
		},
		{
			name: "passcode expired",
			mockedRes: mockedRes{
				modelCallErr: models.ExpiredPasscodeErr,
				jsonCode:     http.StatusBadRequest},
			wantErr: models.ExpiredPasscodeErr,
		},
		{
			name: "passcode mismatches",
			mockedRes: mockedRes{
				modelCallErr: models.InvalidDataErr,
				jsonCode:     http.StatusBadRequest},
			wantErr: models.InvalidDataErr,
		},
		{
			name: "passcode mismatches",
			mockedRes: mockedRes{
				modelCallErr:  err,
				noContentCode: http.StatusInternalServerError},
			wantErr: err,
		},
		{
			name: "success",
			mockedRes: mockedRes{
				noContentCode: http.StatusCreated},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanSetup := setup(tt.mockedRes)
			defer cleanSetup()
			err := user.NewUserController(userModelMock).ConfirmSignup(ctxMock)
			if err != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			assert.Nil(t, err)
		})
	}
}
