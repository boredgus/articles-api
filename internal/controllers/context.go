package controllers

import (
	"net/http"
	"net/url"
)

// success
// swagger:response successResp200
// nolint:unused
type successResp200 struct{}

type InfoResponse struct {
	Message string `json:"message"`
}

// swagger:response respWithMessage
// nolint:unused
type updateUserResp200 struct {
	// success
	// in: body
	Body InfoResponse
}

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

// invalid data provided
// swagger:response invalidData400
// nolint:unused
type invalidData400 struct {
	// in: body
	body ErrorBody
}

// unauthorized
// swagger:response unauthorizedResp401
// nolint:unused
type unauthorizedResp401 struct {
	// in: body
	Body ErrorBody
}

// user does not have enough rights to perform action
// swagger:response forbiddenResp403
// nolint:unused
type forbiddenResp403 struct {
	// in: body
	// required: true
	Body ErrorBody
}

// there is no article with such id
// swagger:response articleNotFound404
// nolint:unused
type articleNotFound404 struct {
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
