package user

import (
	"a-article/internal/auth"
	cntr "a-article/internal/controllers"
	"a-article/internal/models"
	"errors"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

type updateRolePayload struct {
	// role to set for specified user
	// required: true
	// enum: user,moderator
	Role string `json:"role" form:"role"`
}

// swagger:parameters update_user_role
// nolint:unused
type updateRoleParams struct {
	// user identifier
	// in: path
	// required: true
	UserOId string `json:"user_id"`
	// in: body
	// required: true
	Payload updateRolePayload `json:"payload"`
}

// invalid role data
// swagger:response updateUserResp400
// nolint:unused
type updateUserResp400 struct {
	// in: body
	Body cntr.ErrorBody
}

// swagger:route PATCH /users/{user_id}/role users update_user_role
// updates user role
// ---
// Requires admin role. It is prohibited to update user with `admin` role.
//
// security:
//   - jwt:
//
// responses:
//
//	200: respWithMessage
//	400: updateUserResp400
//	401: unauthorizedResp401
//	403: forbiddenResp403
//	404: userNotFound
//	500: commonError
func (u User) UpdateRole(ctx cntr.Context) error {
	var payload updateRolePayload
	err := ctx.Bind(&payload)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, cntr.ErrorBody{Error: "failed to bind data"})
	}
	err = u.userModel.UpdateRole(
		ctx.Get("user").(*jwt.Token).Claims.(*auth.JWTClaims).Role,
		ctx.PathParam("user_id"),
		payload.Role,
	)
	if errors.Is(err, models.InvalidUserDataErr) {
		e := ctx.JSON(http.StatusBadRequest, cntr.ErrorBody{Error: err.Error()})
		return fmt.Errorf("%v: %w", e, err)
	}
	if errors.Is(err, models.UserNotFoundErr) {
		e := ctx.JSON(http.StatusNotFound, cntr.ErrorBody{Error: err.Error()})
		return fmt.Errorf("%v: %w", e, err)
	}
	if errors.Is(err, models.NotEnoughRightsErr) {
		e := ctx.JSON(http.StatusForbidden, cntr.ErrorBody{Error: err.Error()})
		return fmt.Errorf("%v: %w", e, err)
	}
	if err != nil {
		return fmt.Errorf("%v: %w", ctx.NoContent(http.StatusInternalServerError), err)
	}
	return ctx.JSON(http.StatusOK, cntr.InfoResponse{Message: "successfuly updated role to " + payload.Role})
}
