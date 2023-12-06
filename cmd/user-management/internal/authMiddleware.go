package internal

import (
	"net/http"
	"user-management/internal/controllers"
	"user-management/internal/domain"
	"user-management/internal/models"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func authMiddleware(user models.UserModel) middleware.BasicAuthValidator {
	return func(username, password string, c echo.Context) (bool, error) {
		_, _, err := user.Authorize(domain.NewUser(username, password))
		if err != nil {
			return false, c.JSON(http.StatusUnauthorized, controllers.ErrorBody{Error: "user is unauthorized"})
		}
		c.Request().Header.Set("Username", username)
		return true, nil
	}
}
