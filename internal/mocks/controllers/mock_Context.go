// Code generated by mockery v2.38.0. DO NOT EDIT.

package mocks

import (
	http "net/http"
	url "net/url"

	mock "github.com/stretchr/testify/mock"
)

// Context is an autogenerated mock type for the Context type
type Context struct {
	mock.Mock
}

type Context_Expecter struct {
	mock *mock.Mock
}

func (_m *Context) EXPECT() *Context_Expecter {
	return &Context_Expecter{mock: &_m.Mock}
}

// Bind provides a mock function with given fields: i
func (_m *Context) Bind(i interface{}) error {
	ret := _m.Called(i)

	if len(ret) == 0 {
		panic("no return value specified for Bind")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}) error); ok {
		r0 = rf(i)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Context_Bind_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Bind'
type Context_Bind_Call struct {
	*mock.Call
}

// Bind is a helper method to define mock.On call
//   - i interface{}
func (_e *Context_Expecter) Bind(i interface{}) *Context_Bind_Call {
	return &Context_Bind_Call{Call: _e.mock.On("Bind", i)}
}

func (_c *Context_Bind_Call) Run(run func(i interface{})) *Context_Bind_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(interface{}))
	})
	return _c
}

func (_c *Context_Bind_Call) Return(_a0 error) *Context_Bind_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Context_Bind_Call) RunAndReturn(run func(interface{}) error) *Context_Bind_Call {
	_c.Call.Return(run)
	return _c
}

// FormParams provides a mock function with given fields:
func (_m *Context) FormParams() (url.Values, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for FormParams")
	}

	var r0 url.Values
	var r1 error
	if rf, ok := ret.Get(0).(func() (url.Values, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() url.Values); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(url.Values)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Context_FormParams_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FormParams'
type Context_FormParams_Call struct {
	*mock.Call
}

// FormParams is a helper method to define mock.On call
func (_e *Context_Expecter) FormParams() *Context_FormParams_Call {
	return &Context_FormParams_Call{Call: _e.mock.On("FormParams")}
}

func (_c *Context_FormParams_Call) Run(run func()) *Context_FormParams_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Context_FormParams_Call) Return(_a0 url.Values, _a1 error) *Context_FormParams_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Context_FormParams_Call) RunAndReturn(run func() (url.Values, error)) *Context_FormParams_Call {
	_c.Call.Return(run)
	return _c
}

// Get provides a mock function with given fields: key
func (_m *Context) Get(key string) interface{} {
	ret := _m.Called(key)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(string) interface{}); ok {
		r0 = rf(key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	return r0
}

// Context_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type Context_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - key string
func (_e *Context_Expecter) Get(key interface{}) *Context_Get_Call {
	return &Context_Get_Call{Call: _e.mock.On("Get", key)}
}

func (_c *Context_Get_Call) Run(run func(key string)) *Context_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *Context_Get_Call) Return(_a0 interface{}) *Context_Get_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Context_Get_Call) RunAndReturn(run func(string) interface{}) *Context_Get_Call {
	_c.Call.Return(run)
	return _c
}

// JSON provides a mock function with given fields: code, i
func (_m *Context) JSON(code int, i interface{}) error {
	ret := _m.Called(code, i)

	if len(ret) == 0 {
		panic("no return value specified for JSON")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(int, interface{}) error); ok {
		r0 = rf(code, i)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Context_JSON_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'JSON'
type Context_JSON_Call struct {
	*mock.Call
}

// JSON is a helper method to define mock.On call
//   - code int
//   - i interface{}
func (_e *Context_Expecter) JSON(code interface{}, i interface{}) *Context_JSON_Call {
	return &Context_JSON_Call{Call: _e.mock.On("JSON", code, i)}
}

func (_c *Context_JSON_Call) Run(run func(code int, i interface{})) *Context_JSON_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int), args[1].(interface{}))
	})
	return _c
}

func (_c *Context_JSON_Call) Return(_a0 error) *Context_JSON_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Context_JSON_Call) RunAndReturn(run func(int, interface{}) error) *Context_JSON_Call {
	_c.Call.Return(run)
	return _c
}

// NoContent provides a mock function with given fields: code
func (_m *Context) NoContent(code int) error {
	ret := _m.Called(code)

	if len(ret) == 0 {
		panic("no return value specified for NoContent")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(int) error); ok {
		r0 = rf(code)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Context_NoContent_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'NoContent'
type Context_NoContent_Call struct {
	*mock.Call
}

// NoContent is a helper method to define mock.On call
//   - code int
func (_e *Context_Expecter) NoContent(code interface{}) *Context_NoContent_Call {
	return &Context_NoContent_Call{Call: _e.mock.On("NoContent", code)}
}

func (_c *Context_NoContent_Call) Run(run func(code int)) *Context_NoContent_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int))
	})
	return _c
}

func (_c *Context_NoContent_Call) Return(_a0 error) *Context_NoContent_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Context_NoContent_Call) RunAndReturn(run func(int) error) *Context_NoContent_Call {
	_c.Call.Return(run)
	return _c
}

// PathParam provides a mock function with given fields: name
func (_m *Context) PathParam(name string) string {
	ret := _m.Called(name)

	if len(ret) == 0 {
		panic("no return value specified for PathParam")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(name)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Context_PathParam_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'PathParam'
type Context_PathParam_Call struct {
	*mock.Call
}

// PathParam is a helper method to define mock.On call
//   - name string
func (_e *Context_Expecter) PathParam(name interface{}) *Context_PathParam_Call {
	return &Context_PathParam_Call{Call: _e.mock.On("PathParam", name)}
}

func (_c *Context_PathParam_Call) Run(run func(name string)) *Context_PathParam_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *Context_PathParam_Call) Return(_a0 string) *Context_PathParam_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Context_PathParam_Call) RunAndReturn(run func(string) string) *Context_PathParam_Call {
	_c.Call.Return(run)
	return _c
}

// QueryParams provides a mock function with given fields:
func (_m *Context) QueryParams() url.Values {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for QueryParams")
	}

	var r0 url.Values
	if rf, ok := ret.Get(0).(func() url.Values); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(url.Values)
		}
	}

	return r0
}

// Context_QueryParams_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'QueryParams'
type Context_QueryParams_Call struct {
	*mock.Call
}

// QueryParams is a helper method to define mock.On call
func (_e *Context_Expecter) QueryParams() *Context_QueryParams_Call {
	return &Context_QueryParams_Call{Call: _e.mock.On("QueryParams")}
}

func (_c *Context_QueryParams_Call) Run(run func()) *Context_QueryParams_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Context_QueryParams_Call) Return(_a0 url.Values) *Context_QueryParams_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Context_QueryParams_Call) RunAndReturn(run func() url.Values) *Context_QueryParams_Call {
	_c.Call.Return(run)
	return _c
}

// Request provides a mock function with given fields:
func (_m *Context) Request() *http.Request {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Request")
	}

	var r0 *http.Request
	if rf, ok := ret.Get(0).(func() *http.Request); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*http.Request)
		}
	}

	return r0
}

// Context_Request_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Request'
type Context_Request_Call struct {
	*mock.Call
}

// Request is a helper method to define mock.On call
func (_e *Context_Expecter) Request() *Context_Request_Call {
	return &Context_Request_Call{Call: _e.mock.On("Request")}
}

func (_c *Context_Request_Call) Run(run func()) *Context_Request_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Context_Request_Call) Return(_a0 *http.Request) *Context_Request_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Context_Request_Call) RunAndReturn(run func() *http.Request) *Context_Request_Call {
	_c.Call.Return(run)
	return _c
}

// NewContext creates a new instance of Context. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewContext(t interface {
	mock.TestingT
	Cleanup(func())
}) *Context {
	mock := &Context{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
