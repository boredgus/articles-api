package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"user-management/internal/domain"
	"user-management/internal/models"
)

type LoginController interface {
	Register(ctx Context) error
	Authorize(ctx Context) error
}

func NewLoginController(userModel models.UserModel) LoginController {
	return Login{userModel: userModel}
}

type Login struct {
	userModel models.UserModel
}

type AuthBody struct {
	Token  string `json:"token"`
	UserId string `json:"user_id"`
}

func (c Login) Register(ctx Context) error {
	var user domain.User
	err := ctx.Bind(&user)
	if err != nil {
		e := ctx.JSON(http.StatusBadRequest, ErrorBody{Error: "username and password are required"})
		return fmt.Errorf("%v: %w", e, err)
	}

	err = c.userModel.Create(user)
	if errors.Is(err, models.UsernameDuplicationErr) {
		e := ctx.JSON(http.StatusConflict, ErrorBody{Error: err.Error()})
		return fmt.Errorf("%v: %w", e, err)
	}
	if errors.Is(err, models.InvalidAuthParameterErr) {
		e := ctx.JSON(http.StatusBadRequest, ErrorBody{Error: err.Error()})
		return fmt.Errorf("%v: %w", e, err)
	}
	if err != nil {
		e := ctx.NoContent(http.StatusInternalServerError)
		return fmt.Errorf("%v: %w", e, err)
	}

	return ctx.NoContent(http.StatusCreated)
}

func (c Login) Authorize(ctx Context) error {
	var user domain.User
	err := ctx.Bind(&user)
	if err != nil {
		e := ctx.JSON(http.StatusUnauthorized, ErrorBody{Error: "username and password are required"})
		return fmt.Errorf("%v: %w", e, err)
	}

	userId, token, err := c.userModel.Authorize(user)
	if errors.Is(err, models.InvalidAuthParameterErr) {
		e := ctx.JSON(http.StatusUnauthorized, ErrorBody{Error: err.Error()})
		return fmt.Errorf("%v: %w", e, err)
	}
	if err != nil {
		e := ctx.NoContent(http.StatusInternalServerError)
		return fmt.Errorf("%v: %w", e, err)
	}

	return ctx.JSON(http.StatusOK, AuthBody{Token: token, UserId: userId})
}
