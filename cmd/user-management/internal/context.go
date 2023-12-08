package internal

import (
	"net/http"
	"net/url"
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

func (c EchoContext) QueryParams() url.Values {
	return c.echo.QueryParams()
}
func (c EchoContext) FormParams() (url.Values, error) {
	return c.echo.FormParams()
}
func (c EchoContext) Request() *http.Request {
	return c.echo.Request()
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
