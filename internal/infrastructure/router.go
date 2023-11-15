package infrastructure

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type t struct {
	Message string
}

func GetRouter() *echo.Echo {
	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format:           "${time_custom} ${method} ${uri} ${status} error:${error}\n",
		CustomTimeFormat: "2006-01-02 15:04:05",
	}))
	e.GET("/", func(c echo.Context) error {
		c.JSON(http.StatusOK, struct {
			Message string `json:"message"`
		}{Message: "message"})
		return nil
	})

	return e
}
