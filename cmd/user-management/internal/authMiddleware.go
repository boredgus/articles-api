package internal

import (
	"net/http"
	"user-management/internal/auth"
	"user-management/internal/controllers"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func jwtAuthMiddleware() echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		SigningMethod: "HS256",
		SigningKey:    auth.JWTSecretKey,
		ContextKey:    "user",
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(auth.JWTClaims)
		},
		ErrorHandler: func(c echo.Context, err error) error {
			return c.JSON(http.StatusUnauthorized, controllers.ErrorBody{Error: err.Error()})
		},
	})
}
