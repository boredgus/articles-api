package controllers

// error entity
// swagger:model
type ErrorBody struct {
	// error details
	// required: true
	Error string `json:"error"`
}

// internal server error
// swagger:response commonError
type Error struct {
	// in: body
	Body ErrorBody
}

type Context interface {
	PathParam(name string) string
	Bind(i interface{}) error

	NoContent(code int) error
	JSON(code int, i interface{}) error
}
