package controllers

import (
	"fmt"
	"net/http"
	"testing"
	"user-management/internal/domain"
	"user-management/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type loginModelMock struct {
	mock.Mock
}

func (l *loginModelMock) Create(user domain.User) error {
	args := l.Called(user)
	res := args.Get(0)
	if res != nil {
		return res.(error)
	}
	return nil
}
func (l *loginModelMock) Authorize(user domain.User) (string, string, error) {
	args := l.Called(user)
	userId, token, err := args.Get(0).(string), args.Get(1).(string), args.Get(2)
	if err != nil {
		return userId, token, err.(error)
	}
	return userId, token, nil
}

func TestUserController(t *testing.T) {
	model := new(loginModelMock)
	context := new(contextMock)
	serverError := fmt.Errorf("server err")
	bindErr := fmt.Errorf("bind err")

	t.Run("registration fails: user binding failed", func(t *testing.T) {
		bindMock := context.On("Bind", mock.Anything).Return(bindErr)
		jsonMock := context.On("JSON", http.StatusBadRequest, mock.Anything).Return(nil)

		err := NewLoginController(model).Register(context)
		assert.ErrorIs(t, err, bindErr)

		bindMock.Unset()
		jsonMock.Unset()
	})

	t.Run("registration fails if user with such username already exists", func(t *testing.T) {
		modelMock := model.On("Create", mock.Anything).Return(models.UsernameDuplicationErr)
		bindMock := context.On("Bind", mock.Anything).Return(nil)
		jsonMock := context.On("JSON", http.StatusConflict, mock.Anything).Return(nil)

		err := NewLoginController(model).Register(context)
		assert.ErrorIs(t, err, models.UsernameDuplicationErr)

		modelMock.Unset()
		bindMock.Unset()
		jsonMock.Unset()
	})

	t.Run("registration fails if user credentials are invalid", func(t *testing.T) {
		modelMock := model.On("Create", mock.Anything).Return(models.InvalidAuthParameterErr)
		bindMock := context.On("Bind", mock.Anything).Return(nil)
		jsonMock := context.On("JSON", http.StatusBadRequest, mock.Anything).Return(nil)

		err := NewLoginController(model).Register(context)
		assert.ErrorIs(t, err, models.InvalidAuthParameterErr)

		modelMock.Unset()
		bindMock.Unset()
		jsonMock.Unset()
	})

	t.Run("registration fails because of internal error", func(t *testing.T) {
		modelMock := model.On("Create", mock.Anything).Return(serverError)
		bindMock := context.On("Bind", mock.Anything).Return(nil)
		jsonMock := context.On("NoContent", http.StatusInternalServerError).Return(nil)

		err := NewLoginController(model).Register(context)
		assert.ErrorIs(t, err, serverError)

		modelMock.Unset()
		bindMock.Unset()
		jsonMock.Unset()
	})

	t.Run("success registration", func(t *testing.T) {
		modelMock := model.On("Create", mock.Anything).Return(nil)
		bindMock := context.On("Bind", mock.Anything).Return(nil)
		jsonMock := context.On("NoContent", http.StatusCreated).Return(nil)

		err := NewLoginController(model).Register(context)
		assert.Nil(t, err)

		modelMock.Unset()
		bindMock.Unset()
		jsonMock.Unset()
	})

	t.Run("authorization fails: user binding failed", func(t *testing.T) {
		bindMock := context.On("Bind", mock.Anything).Return(bindErr)
		jsonMock := context.On("JSON", http.StatusUnauthorized, mock.Anything).Return(nil)

		err := NewLoginController(model).Authorize(context)
		assert.ErrorIs(t, err, bindErr)

		bindMock.Unset()
		jsonMock.Unset()
	})

	t.Run("authorization fails if user credentials are invalid", func(t *testing.T) {
		bindMock := context.On("Bind", mock.Anything).Return(nil)
		modelMock := model.On("Authorize", mock.Anything).Return("", "", models.InvalidAuthParameterErr)
		jsonMock := context.On("JSON", http.StatusUnauthorized, mock.Anything).Return(nil)

		err := NewLoginController(model).Authorize(context)
		assert.ErrorIs(t, err, models.InvalidAuthParameterErr)

		bindMock.Unset()
		modelMock.Unset()
		jsonMock.Unset()
	})

	t.Run("authorization fails because of internal error", func(t *testing.T) {
		bindMock := context.On("Bind", mock.Anything).Return(nil)
		modelMock := model.On("Authorize", mock.Anything).Return("", "", serverError)
		jsonMock := context.On("NoContent", http.StatusInternalServerError, mock.Anything).Return(nil)

		err := NewLoginController(model).Authorize(context)
		assert.ErrorIs(t, err, serverError)

		bindMock.Unset()
		modelMock.Unset()
		jsonMock.Unset()
	})

	t.Run("successful authorization", func(t *testing.T) {
		bindMock := context.On("Bind", mock.Anything).Return(nil)
		userId, token := "user-id", "token"
		modelMock := model.On("Authorize", mock.Anything).Return(userId, token, nil)
		jsonMock := context.On("JSON", http.StatusOK, AuthBody{Token: token, UserId: userId}).Return(nil)

		err := NewLoginController(model).Authorize(context)
		assert.Nil(t, err)

		bindMock.Unset()
		modelMock.Unset()
		jsonMock.Unset()
	})
}
