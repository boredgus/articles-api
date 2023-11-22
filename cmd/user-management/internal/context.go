package internal

import (
	"user-management/internal/controllers"

	"github.com/labstack/echo/v4"
)

func NewContext(e echo.Context) controllers.Context {
	return EchoContext{
		echo: e,
	}
}

type EchoContext struct {
	echo echo.Context
}

func (c EchoContext) NoContent(code int) error {
	return c.echo.NoContent(code)
}

func (c EchoContext) JSON(code int, i interface{}) error {
	return c.echo.JSON(code, i)
}

func (c EchoContext) PathParam(name string) string {
	return c.echo.Param(name)
}

func (c EchoContext) Bind(i interface{}) error {
	return c.echo.Bind(i)
}
