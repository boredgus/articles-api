package user

import (
	"errors"
	"fmt"
	"net/http"
	cntrl "user-management/internal/controllers"
	"user-management/internal/domain"
	"user-management/internal/models"
)

// swagger:parameters registration
// nolint:unused
type registerParams struct {
	// secret key to create user with protected role
	// in: query
	APIKey string `json:"api_key"`
	// in: body
	// required: true
	Body domain.User `json:"user"`
}

// user created
// swagger:response registerResp201
// nolint:unused
type registerResp201 struct{}

// user with such username already exists
// swagger:response registerResp409
// nolint:unused
type registerResp409 struct {
	// in: body
	body cntrl.ErrorBody
}

// swagger:route POST /register login registration
// creates new user
// ---
// Checks whether user with provided username exists and validates provided password
// responses:
//
//		201: registerResp201
//	 	409: registerResp409
//		500: commonError
func (c Login) Register(ctx cntrl.Context) error {
	var user domain.User
	err := ctx.Bind(&user)
	if err != nil {
		e := ctx.JSON(http.StatusBadRequest, cntrl.ErrorBody{Error: "username and password are required"})
		return fmt.Errorf("%v: %w", e, err)
	}

	err = c.userModel.Create(user, ctx.QueryParams().Get("api_key"))
	if errors.Is(err, models.UsernameDuplicationErr) {
		e := ctx.JSON(http.StatusConflict, cntrl.ErrorBody{Error: err.Error()})
		return fmt.Errorf("%v: %w", e, err)
	}
	if errors.Is(err, models.InvalidUserErr) || errors.Is(err, models.InvalidAPIKeyErr) {
		e := ctx.JSON(http.StatusBadRequest, cntrl.ErrorBody{Error: err.Error()})
		return fmt.Errorf("%v: %w", e, err)
	}
	if err != nil {
		e := ctx.NoContent(http.StatusInternalServerError)
		return fmt.Errorf("%v: %w", e, err)
	}

	return ctx.NoContent(http.StatusCreated)
}
