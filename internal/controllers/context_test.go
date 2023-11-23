package controllers

import "github.com/stretchr/testify/mock"

type contextMock struct {
	mock.Mock
}

func (c *contextMock) PathParam(name string) string {
	args := c.Called(name)
	return args.Get(0).(string)
}
func (c *contextMock) Bind(i interface{}) error {
	args := c.Called(i)
	res := args.Get(0)
	if res != nil {
		return res.(error)
	}
	return nil
}
func (c *contextMock) NoContent(code int) error {
	args := c.Called(code)
	res := args.Get(0)
	if res != nil {
		return res.(error)
	}
	return nil
}
func (c *contextMock) JSON(code int, i interface{}) error {
	args := c.Called(code, i)
	res := args.Get(0)
	if res != nil {
		return res.(error)
	}
	return nil
}
