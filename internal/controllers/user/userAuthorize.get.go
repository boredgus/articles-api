package user

import (
	cntrl "a-article/internal/controllers"
	"a-article/internal/domain"
	"a-article/internal/models"
	"errors"
	"fmt"
	"net/http"
)

// swagger:parameters authorization
// nolint:unused
type authParams struct {
	// in: body
	// required: true
	Body domain.User `json:"user"`
}

// access token and user id
// swagger:model
type authResult struct {
	// access token
	// required: true
	Token string `json:"token"`
}

// successfully authorized
// swagger:response authResp200
// nolint:unused
type authResp200 struct {
	// in: body
	body authResult
}

// username or password is invalid
// swagger:response authResp401
// nolint:unused
type authResp401 struct {
	// in: body
	body cntrl.ErrorBody
}

// swagger:route POST /authorize auth authorization
// authorizes user
// ---
// Checks whether user with such username exists and compares his password with given one.
//
// responses:
//
//	 	200: authResp200
//		401: authResp401
//		500: commonError
func (c User) Authorize(ctx cntrl.Context) error {
	var user domain.User
	err := ctx.Bind(&user)
	if err != nil {
		e := ctx.JSON(http.StatusUnauthorized, cntrl.ErrorBody{Error: "username and password are required"})
		return fmt.Errorf("%v: %w", e, err)
	}

	token, err := c.userModel.Authorize(user.Username, user.Password)
	if errors.Is(err, models.InvalidAuthParameterErr) {
		e := ctx.JSON(http.StatusUnauthorized, cntrl.ErrorBody{Error: err.Error()})
		return fmt.Errorf("%v: %w", e, err)
	}
	if err != nil {
		e := ctx.NoContent(http.StatusInternalServerError)
		return fmt.Errorf("%v: %w", e, err)
	}

	return ctx.JSON(http.StatusOK, authResult{Token: token})
}
