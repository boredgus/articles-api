package user

import (
	cntrl "a-article/internal/controllers"
	"a-article/internal/domain"
	"a-article/internal/models"
	"errors"
	"fmt"
	"net/http"
)

// swagger:parameters registration
// nolint:unused
type registerParams struct {
	// in: body
	// required: true
	Body domain.User `json:"user"`
}

// passcode was sent to provided email to configm signup
// swagger:response registerResp200
// nolint:unused
type registerResp200 struct{}

// user with such username already exists
// swagger:response registerResp409
// nolint:unused
type registerResp409 struct {
	// in: body
	body cntrl.ErrorBody
}

// swagger:route POST /register login registration
// starts signup of new user
// ---
// Checks whether user with provided username exists and validates provided password
// responses:
//
//		200: registerResp200
//	 	409: registerResp409
//		500: commonError
func (c User) Register(ctx cntrl.Context) error {
	// TODO: use separate input/output types instead of entities in controllers
	var user domain.User
	err := ctx.Bind(&user)
	if err != nil {
		e := ctx.JSON(http.StatusBadRequest, cntrl.ErrorBody{Error: "username and password are required"})
		return fmt.Errorf("%v: %w", e, err)
	}

	err = c.userModel.RequestSignup(user)
	if errors.Is(err, models.UsernameDuplicationErr) {
		e := ctx.JSON(http.StatusConflict, cntrl.ErrorBody{Error: err.Error()})
		return fmt.Errorf("%v: %w", e, err)
	}
	if errors.Is(err, models.InvalidDataErr) {
		e := ctx.JSON(http.StatusBadRequest, cntrl.ErrorBody{Error: err.Error()})
		return fmt.Errorf("%v: %w", e, err)
	}
	if err != nil {
		e := ctx.NoContent(http.StatusInternalServerError)
		return fmt.Errorf("%v: %w", e, err)
	}

	return ctx.JSON(http.StatusOK, cntrl.InfoResponse{Message: "passcode was sent to " + user.Username})
}
