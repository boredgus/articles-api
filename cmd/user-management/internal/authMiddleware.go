package internal

import (
	"user-management/internal/domain"
	"user-management/internal/models"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func authMiddleware(model models.UserModel) middleware.BasicAuthValidator {
	return func(username, password string, c echo.Context) (bool, error) {
		_, _, err := model.Authorize(domain.NewUser(username, password))
		if err != nil {
			return false, nil
		}
		c.Request().Header.Set("Username", username)
		return true, nil
	}
}
