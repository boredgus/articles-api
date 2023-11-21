package controllers

import (
	"fmt"
	"net/http"
	"user-management/internal/auth"
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

	if len(user.Username) == 0 || len(user.Password) == 0 {
		ctx.JSON(http.StatusBadRequest, ErrorBody{Error: "username and password cannot be empty"})
		return fmt.Errorf("username or password is empty")
	}

	user.Password = auth.NewPassword().Hash(user.Password)
	err = c.userModel.Create(user)
	if err == models.UsernameDuplicationErr {
		ctx.JSON(http.StatusConflict, ErrorBody{Error: "user with such username already exists"})
		return err
	}
	if err != nil {
		ctx.NoContent(http.StatusInternalServerError)
		return err
	}

	ctx.NoContent(http.StatusCreated)
	return nil
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
	if err == models.InvalidAuthParameterErr {
		ctx.JSON(http.StatusUnauthorized, ErrorBody{Error: err.Error()})
		return err
	}
	if err != nil {
		ctx.NoContent(http.StatusInternalServerError)
		return err
	}

	ctx.JSON(http.StatusOK, AuthBody{Token: token, UserId: userId})
	return nil
}
