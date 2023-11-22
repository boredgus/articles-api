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
		ctx.JSON(http.StatusBadRequest, ErrorBody{Error: "username and password are required"})
		return err
	}

	// if err = user.Validate(); err != nil {
	// 	ctx.JSON(http.StatusBadRequest, ErrorBody{Error: err.Error()})
	// 	return err
	// }

	// hashedPswd, err := auth.NewPassword().Hash(user.Password)
	// if err != nil {

	// 	ctx.NoContent(http.StatusInternalServerError)
	// 	return err
	// }

	// user.Password = hashedPswd
	err = c.userModel.Create(user)
	if errors.Is(err, models.UsernameDuplicationErr) {
		ctx.JSON(http.StatusConflict, ErrorBody{Error: err.Error()})
		return err
	}
	if errors.Is(err, models.InvalidAuthParameterErr) {
		ctx.JSON(http.StatusBadRequest, ErrorBody{Error: err.Error()})
		return err
	}
	if err != nil {
		ctx.NoContent(http.StatusInternalServerError)
		return err
	}

	return ctx.NoContent(http.StatusCreated)
}

func (c Login) Authorize(ctx Context) error {
	var user domain.User
	err := ctx.Bind(&user)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, ErrorBody{Error: "username and password are required"})
		return err
	}

	if len(user.Username) == 0 || len(user.Password) == 0 {
		ctx.JSON(http.StatusUnauthorized, ErrorBody{Error: "username and password cannot be empty"})
		return fmt.Errorf("username or password is empty")
	}

	userId, token, err := c.userModel.Authorize(user)
	if errors.Is(err, models.InvalidAuthParameterErr) {
		ctx.JSON(http.StatusUnauthorized, ErrorBody{Error: err.Error()})
		return err
	}
	if err != nil {
		ctx.NoContent(http.StatusInternalServerError)
		return err
	}

	return ctx.JSON(http.StatusOK, AuthBody{Token: token, UserId: userId})
}
