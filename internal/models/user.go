package models

import (
	"errors"
	"fmt"
	"user-management/internal/auth"
	"user-management/internal/domain"
	"user-management/internal/models/repo"

	"github.com/google/uuid"
)

type UserModel interface {
	Create(user domain.User) error
	Authorize(user domain.User) (string, string, error)
	Exists(oid, password string) error
}

var InvalidUserErr = errors.New("invalid user data")
var InvalidAPIKeyErr = errors.New("invalid api key")
var InvalidAuthParameterErr = errors.New("username or password is invalid")
var UsernameDuplicationErr = errors.New("user with such username already exists")
var UserNotFoundErr = errors.New("user not found")

func NewUserModel(repo repo.UserRepository) UserModel {
	return user{repo: repo, token: auth.NewToken(), pswd: auth.NewPassword()}
}

type user struct {
	repo  repo.UserRepository
	token auth.Token
	pswd  auth.Password
}

func (u user) Create(user domain.User) error {
	if err := user.Validate(); err != nil {
		return fmt.Errorf("%w: %w", InvalidUserErr, err)
	}
	hashedPswd, err := u.pswd.Hash(user.Password)
	if err != nil {
		return err
	}
	return u.repo.Create(repo.User{
		OId:      uuid.New().String(),
		Username: user.Username,
		Password: hashedPswd,
		Role:     repo.DefaultUserRole,
	})
}

func (u user) Authorize(user domain.User) (userId string, token string, err error) {
	userFromDB, err := u.repo.Get(user.Username)
	if err != nil {
		return "", "", InvalidAuthParameterErr
	}
	if !u.pswd.Compare(userFromDB.Password, user.Password) {
		return "", "", InvalidAuthParameterErr
	}

	token, err = u.token.Generate(user)
	return userFromDB.OId, token, err
}

func (u user) Exists(oid, password string) error {
	userFromDB, err := u.repo.GetByOId(oid)
	if err != nil {
		return err
	}
	if !u.pswd.Compare(userFromDB.Password, password) {
		return InvalidAuthParameterErr
	}
	return nil
}
