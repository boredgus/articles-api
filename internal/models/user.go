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
	Authorize(username, password string) (string, error)
	Exists(oid, password string) error
}

var InvalidUserErr = errors.New("invalid user data")
var InvalidAPIKeyErr = errors.New("invalid api key")
var InvalidAuthParameterErr = errors.New("username or password is invalid")
var UsernameDuplicationErr = errors.New("user with such username already exists")
var UserNotFoundErr = errors.New("user not found")

func NewUserModel(repo repo.UserRepository) UserModel {
	return user{repo: repo, token: auth.NewJWT(), pswd: auth.NewPassword()}
}

type user struct {
	repo  repo.UserRepository
	token auth.Token[auth.JWTPayload]
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
		Role:     domain.DefaultUserRole,
	})
}

func (u user) Authorize(username, password string) (token string, err error) {
	userFromDB, err := u.repo.Get(username)
	if err != nil {
		return "", InvalidAuthParameterErr
	}
	if !u.pswd.Compare(userFromDB.Password, password) {
		return "", InvalidAuthParameterErr
	}
	return u.token.Generate(auth.JWTPayload{
		Username: userFromDB.Username,
		UserOId:  userFromDB.OId,
		Role:     domain.UserRoles[userFromDB.Role],
	})
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
