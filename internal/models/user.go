package models

import (
	"errors"
	"user-management/internal/auth"
	"user-management/internal/domain"
)

type UserModel interface {
	Create(user domain.User) error
	Authorize(user domain.User) (string, error)
}

type UserRepository interface {
	Create(user domain.User) error
	Get(username string) (domain.User, error)
}

var InvalidAuthParameter = errors.New("username or password is invalid")

func NewUserModel(repo UserRepository) UserModel {
	return User{repo: repo, token: auth.NewToken()}
}

type User struct {
	repo  UserRepository
	token auth.Token
}

func (u User) Create(user domain.User) error {
	return u.repo.Create(user)
}

func (u User) Authorize(user domain.User) (string, error) {
	userFromDB, err := u.repo.Get(user.Username)
	if err != nil {
		return "", InvalidAuthParameter
	}
	arePswdsEqual := auth.NewPassword().Compare(userFromDB.Password, user.Password)
	if !arePswdsEqual {
		return "", InvalidAuthParameter
	}

	return u.token.Generate(user)
}
