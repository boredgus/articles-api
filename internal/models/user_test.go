package models

import (
	"testing"
	"user-management/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type userRepoMock struct {
	mock.Mock
}

func NewUserRepoMock() userRepoMock {
	return userRepoMock{}
}
func (r *userRepoMock) Create(user User) error {
	args := r.Called(user)
	if args.Get(0) != nil {
		return args.Get(0).(error)
	}
	return nil
}
func (r *userRepoMock) Get(username string) (User, error) {
	args := r.Called(username)
	if args.Get(1) != nil {
		return args.Get(0).(User), args.Get(1).(error)
	}
	return args.Get(0).(User), nil
}

func TestUserModel(t *testing.T) {
	r := NewUserRepoMock()
	repo := &r
	validUser := domain.NewUser("username", "PASsword/123")
	validUserToken := "dXNlcm5hbWU6UEFTc3dvcmQvMTIz"
	validUserFromDB := User{
		OId:      "f7873e08-787b-45e9-b22f-82bdf505cca5",
		Username: "username",
		Password: "$2a$10$YCuxL/v4Rn07gP/ggFZIXeIxj6W9BhTaTj1CBDFH0Qysp4ZpI6Pw6",
	}

	t.Run("user creation fails: invalid credentials", func(t *testing.T) {
		err := NewUserModel(repo).Create(domain.NewUser("qw", "ps"))
		assert.ErrorIs(t, err, InvalidAuthParameterErr)
	})

	t.Run("user creation fails: username is duplicated", func(t *testing.T) {
		mock := repo.On("Create", mock.Anything).Return(UsernameDuplicationErr)
		err := NewUserModel(repo).Create(validUser)
		assert.Equal(t, UsernameDuplicationErr, err)
		mock.Unset()
	})

	t.Run("user creation successes", func(t *testing.T) {
		mock := repo.On("Create", mock.Anything).Return(nil)
		err := NewUserModel(repo).Create(validUser)
		assert.Nil(t, err)
		mock.Unset()
	})

	t.Run("user authorization fails: invalid username", func(t *testing.T) {
		mock := repo.On("Get", mock.Anything).Return(User{}, InvalidAuthParameterErr)
		userId, token, err := NewUserModel(repo).Authorize(validUser)
		assert.Equal(t, "", userId, "user id")
		assert.Equal(t, "", token, "token")
		assert.Equal(t, InvalidAuthParameterErr, err, "error")
		mock.Unset()
	})

	t.Run("user authorization fails: invalid password", func(t *testing.T) {
		mock := repo.On("Get", mock.Anything).Return(validUserFromDB, nil)
		userId, token, err := NewUserModel(repo).Authorize(domain.NewUser(validUser.Username, ""))
		assert.Equal(t, "", userId, "user id")
		assert.Equal(t, "", token, "token")
		assert.Equal(t, InvalidAuthParameterErr, err, "error")
		mock.Unset()
	})

	t.Run("success user authorization", func(t *testing.T) {
		mock := repo.On("Get", mock.Anything).Return(validUserFromDB, nil)
		userId, token, err := NewUserModel(repo).Authorize(validUser)
		assert.Equal(t, validUserFromDB.OId, userId, "user id")
		assert.Equal(t, validUserToken, token, "token")
		assert.Nil(t, err, "error")
		mock.Unset()
	})
}
