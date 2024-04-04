package internal

import (
	"a-article/internal/controllers"
	"net/http"
	"net/url"

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

func (c EchoContext) Get(key string) interface{} {
	return c.echo.Get(key)
}
func (c EchoContext) QueryParams() url.Values {
	return c.echo.QueryParams()
}
func (c EchoContext) FormValue(key string) string {
	return c.echo.FormValue(key)
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
