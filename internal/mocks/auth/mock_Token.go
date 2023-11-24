// Code generated by mockery v2.38.0. DO NOT EDIT.

package mocks

import (
	domain "user-management/internal/domain"

	mock "github.com/stretchr/testify/mock"
)

// Token is an autogenerated mock type for the Token type
type Token struct {
	mock.Mock
}

type Token_Expecter struct {
	mock *mock.Mock
}

func (_m *Token) EXPECT() *Token_Expecter {
	return &Token_Expecter{mock: &_m.Mock}
}

// Decode provides a mock function with given fields: token
func (_m *Token) Decode(token string) (domain.User, error) {
	ret := _m.Called(token)

	if len(ret) == 0 {
		panic("no return value specified for Decode")
	}

	var r0 domain.User
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (domain.User, error)); ok {
		return rf(token)
	}
	if rf, ok := ret.Get(0).(func(string) domain.User); ok {
		r0 = rf(token)
	} else {
		r0 = ret.Get(0).(domain.User)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Token_Decode_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Decode'
type Token_Decode_Call struct {
	*mock.Call
}

// Decode is a helper method to define mock.On call
//   - token string
func (_e *Token_Expecter) Decode(token interface{}) *Token_Decode_Call {
	return &Token_Decode_Call{Call: _e.mock.On("Decode", token)}
}

func (_c *Token_Decode_Call) Run(run func(token string)) *Token_Decode_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *Token_Decode_Call) Return(_a0 domain.User, _a1 error) *Token_Decode_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Token_Decode_Call) RunAndReturn(run func(string) (domain.User, error)) *Token_Decode_Call {
	_c.Call.Return(run)
	return _c
}

// Generate provides a mock function with given fields: user
func (_m *Token) Generate(user domain.User) (string, error) {
	ret := _m.Called(user)

	if len(ret) == 0 {
		panic("no return value specified for Generate")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(domain.User) (string, error)); ok {
		return rf(user)
	}
	if rf, ok := ret.Get(0).(func(domain.User) string); ok {
		r0 = rf(user)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(domain.User) error); ok {
		r1 = rf(user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Token_Generate_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Generate'
type Token_Generate_Call struct {
	*mock.Call
}

// Generate is a helper method to define mock.On call
//   - user domain.User
func (_e *Token_Expecter) Generate(user interface{}) *Token_Generate_Call {
	return &Token_Generate_Call{Call: _e.mock.On("Generate", user)}
}

func (_c *Token_Generate_Call) Run(run func(user domain.User)) *Token_Generate_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(domain.User))
	})
	return _c
}

func (_c *Token_Generate_Call) Return(_a0 string, _a1 error) *Token_Generate_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Token_Generate_Call) RunAndReturn(run func(domain.User) (string, error)) *Token_Generate_Call {
	_c.Call.Return(run)
	return _c
}

// NewToken creates a new instance of Token. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewToken(t interface {
	mock.TestingT
	Cleanup(func())
}) *Token {
	mock := &Token{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
