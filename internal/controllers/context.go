package controllers

type ErrorBody struct {
	Error string `json:"error"`
}

type Context interface {
	PathParam(name string) string
	Bind(i interface{}) error

	NoContent(code int) error
	JSON(code int, i interface{}) error
}
