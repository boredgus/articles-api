package models_test

import (
	"testing"
	"user-management/internal/domain"
	"user-management/internal/mocks"
	m "user-management/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserModel(t *testing.T) {
	repo := mocks.NewUserRepository(t)
	validUser := domain.NewUser("username", "PASsword/123")
	validUserToken := "dXNlcm5hbWU6UEFTc3dvcmQvMTIz"
	validUserFromDB := m.User{
		OId:      "f7873e08-787b-45e9-b22f-82bdf505cca5",
		Username: "username",
		Password: "$2a$10$YCuxL/v4Rn07gP/ggFZIXeIxj6W9BhTaTj1CBDFH0Qysp4ZpI6Pw6",
	}

	t.Run("user creation fails: invalid credentials", func(t *testing.T) {
		err := m.NewUserModel(repo).Create(domain.NewUser("qw", "ps"))
		assert.ErrorIs(t, err, m.InvalidAuthParameterErr)
	})

	t.Run("user creation fails: username is duplicated", func(t *testing.T) {
		mock := repo.On("Create", mock.Anything).Return(m.UsernameDuplicationErr)
		err := m.NewUserModel(repo).Create(validUser)
		assert.Equal(t, m.UsernameDuplicationErr, err)
		mock.Unset()
	})

	t.Run("user creation successes", func(t *testing.T) {
		mock := repo.On("Create", mock.Anything).Return(nil)
		err := m.NewUserModel(repo).Create(validUser)
		assert.Nil(t, err)
		mock.Unset()
	})

	t.Run("user authorization fails: invalid username", func(t *testing.T) {
		mock := repo.On("Get", mock.Anything).Return(m.User{}, m.InvalidAuthParameterErr)
		userId, token, err := m.NewUserModel(repo).Authorize(validUser)
		assert.Equal(t, "", userId, "user id")
		assert.Equal(t, "", token, "token")
		assert.Equal(t, m.InvalidAuthParameterErr, err, "error")
		mock.Unset()
	})

	t.Run("user authorization fails: invalid password", func(t *testing.T) {
		mock := repo.On("Get", mock.Anything).Return(validUserFromDB, nil)
		userId, token, err := m.NewUserModel(repo).Authorize(domain.NewUser(validUser.Username, ""))
		assert.Equal(t, "", userId, "user id")
		assert.Equal(t, "", token, "token")
		assert.Equal(t, m.InvalidAuthParameterErr, err, "error")
		mock.Unset()
	})

	t.Run("success user authorization", func(t *testing.T) {
		mock := repo.On("Get", mock.Anything).Return(validUserFromDB, nil)
		userId, token, err := m.NewUserModel(repo).Authorize(validUser)
		assert.Equal(t, validUserFromDB.OId, userId, "user id")
		assert.Equal(t, validUserToken, token, "token")
		assert.Nil(t, err, "error")
		mock.Unset()
	})
}
