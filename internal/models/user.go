package models

import (
	"errors"
	"user-management/internal/auth"
	"user-management/internal/domain"

	"github.com/google/uuid"
)

type User struct {
	OId      string `sql:"o_id"`
	Username string `sql:"username"`
	Password string `sql:"pswd"`
}

type UserModel interface {
	Create(user domain.User) error
	Authorize(user domain.User) (string, string, error)
}

type UserRepository interface {
	Create(user User) error
	Get(username string) (User, error)
}

var InvalidAuthParameterErr = errors.New("username or password is invalid")
var UsernameDuplicationErr = errors.New("user with such username already exists")

func NewUserModel(repo UserRepository) UserModel {
	return user{repo: repo, token: auth.NewToken()}
}

type user struct {
	repo  UserRepository
	token auth.Token
}

func (u user) Create(user domain.User) error {
	return u.repo.Create(User{OId: uuid.New().String(), Username: user.Username, Password: user.Password})
}

func (u user) Authorize(user domain.User) (userId string, token string, err error) {
	userFromDB, err := u.repo.Get(user.Username)
	if err != nil {
		return "", "", InvalidAuthParameterErr
	}
	arePswdsEqual := auth.NewPassword().Compare(userFromDB.Password, user.Password)
	if !arePswdsEqual {
		return "", "", InvalidAuthParameterErr
	}

	token, err = u.token.Generate(user)
	return userFromDB.OId, token, err
}
