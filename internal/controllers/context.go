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
	// in: body
	Body ErrorBody
}

// unauthorized
// swagger:response unauthorizedResp
// nolint:unused
type updateArticleResp401 struct {
	// in: body
	Body ErrorBody
}

type Context interface {
	QueryParams() url.Values
	FormParams() (url.Values, error)
	Request() *http.Request
	PathParam(name string) string
	Bind(i interface{}) error

	NoContent(code int) error
	JSON(code int, i interface{}) error
}
