package controllers

import (
	"net/http"
	"net/url"
)

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
}

// unauthorized
// swagger:response unauthorizedResp401
// nolint:unused
type unauthorizedResp401 struct {
	// in: body
	Body ErrorBody
}

// user does not enough rights to perform action
// swagger:response forbiddenResp403
// nolint:unused
type forbiddenResp403 struct {
	// in: body
	// required: true
	Body ErrorBody
}

// there is no article with such id
// swagger:response notFoundResp404
// nolint:unused
type notFoundResp404 struct {
	// in: body
	// required: true
	Body ErrorBody
}

type Context interface {
	Get(key string) interface{}
	QueryParams() url.Values
	FormParams() (url.Values, error)
	Request() *http.Request
	PathParam(name string) string
	Bind(i interface{}) error

	NoContent(code int) error
	JSON(code int, i interface{}) error
}
