package user

import (
	"errors"
	"fmt"
	"net/http"
	"user-management/internal/auth"
	"user-management/internal/controllers"
	"user-management/internal/models"

	"github.com/golang-jwt/jwt/v5"
)

// swagger:parameters update_user delete_user
// nolint:unused
type updateArticleParams struct {
	// user identifier
	// in: path
	// required: true
	UserId string `json:"user_id"`
}

// user with such id not found
// swagger:response userNotFound
// nolint:unused
type userNotFoundResp struct {
	// in: body
	Body controllers.ErrorBody
}

// swagger:route DELETE /users/{user_id} users delete_user
// deletes user
// ---
// Checks role in token (`admin` role required). It is prohibited to delete user with `admin` role.
//
// security:
//   - jwt:
//
// responses:
//
//	200: successResp200
//	401: unauthorizedResp401
//	403: forbiddenResp403
//	404: userNotFound
//	500: commonError
func (u User) Delete(ctx controllers.Context) error {
	err := u.userModel.Delete(
		ctx.Get("user").(*jwt.Token).Claims.(*auth.JWTClaims).Role,
		ctx.PathParam("user_id"),
	)
	if errors.Is(err, models.UserNotFoundErr) {
		e := ctx.JSON(http.StatusNotFound, controllers.ErrorBody{Error: err.Error()})
		return fmt.Errorf("%v: %w", e, err)
	}
	if errors.Is(err, models.NotEnoughRightsErr) {
		e := ctx.JSON(http.StatusForbidden, controllers.ErrorBody{Error: err.Error()})
		return fmt.Errorf("%v: %w", e, err)
	}
	if err != nil {
		return fmt.Errorf("%v: %w", ctx.NoContent(http.StatusInternalServerError), err)
	}
	return ctx.NoContent(http.StatusOK)
}
