package user

import (
	cntr "a-article/internal/controllers"
	"a-article/internal/models"
	"errors"
	"fmt"
	"net/http"
)

// type confirmParams

// swagger:parameters confirm_signup
// nolint:unused
type confirmSignupParams struct {
	// in: body
	// required: true
	Body struct {
		Username string `json:"username" form:"username"`
		Passcode string `json:"passcode" form:"passcode"`
	} `json:"params"`
}

// passcode is not correct or is expired
// swagger:response confirmSignupResp400
// nolint:unused
type confirmSignupResp400 struct {
	// in: body
	Body cntr.ErrorBody
}

// swagger:route POST /confirm_signup auth confirm_signup
// completes signup
// ---
// Compares given passcode with one sent to email and creates new user.
//
// responses:
//
//		201: successResp200
//		400: confirmSignupResp400
//		404: userNotFound
//	  409: registerResp409
//		500: commonError
func (u User) ConfirmSignup(ctx cntr.Context) error {
	email := ctx.FormValue("username")
	err := u.userModel.ConfirmSignup(email, ctx.FormValue("passcode"))
	if errors.Is(err, models.UsernameDuplicationErr) {
		return fmt.Errorf("%v: %w", ctx.JSON(http.StatusConflict, cntr.ErrorBody{Error: err.Error()}), err)
	}
	if errors.Is(err, models.NotFoundErr) {
		return fmt.Errorf("%w: %w", ctx.JSON(http.StatusNotFound, cntr.ErrorBody{Error: "no signup request found for " + email}), err)
	}
	if errors.Is(err, models.ExpiredPasscodeErr) || errors.Is(err, models.InvalidDataErr) {
		return fmt.Errorf("%w: %w", ctx.JSON(http.StatusBadRequest, cntr.ErrorBody{Error: err.Error()}), err)
	}
	if err != nil {
		return fmt.Errorf("%v: %w", ctx.NoContent(http.StatusInternalServerError), err)
	}
	return ctx.NoContent(http.StatusCreated)
}
