// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// Token is an autogenerated mock type for the Token type
type Token[T interface{}] struct {
	mock.Mock
}

type Token_Expecter[T interface{}] struct {
	mock *mock.Mock
}

func (_m *Token[T]) EXPECT() *Token_Expecter[T] {
	return &Token_Expecter[T]{mock: &_m.Mock}
}

// Decode provides a mock function with given fields: token
func (_m *Token[T]) Decode(token string) (T, error) {
	ret := _m.Called(token)

	if len(ret) == 0 {
		panic("no return value specified for Decode")
	}

	var r0 T
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (T, error)); ok {
		return rf(token)
	}
	if rf, ok := ret.Get(0).(func(string) T); ok {
		r0 = rf(token)
	} else {
		r0 = ret.Get(0).(T)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Token_Decode_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Decode'
type Token_Decode_Call[T interface{}] struct {
	*mock.Call
}

// Decode is a helper method to define mock.On call
//   - token string
func (_e *Token_Expecter[T]) Decode(token interface{}) *Token_Decode_Call[T] {
	return &Token_Decode_Call[T]{Call: _e.mock.On("Decode", token)}
}

func (_c *Token_Decode_Call[T]) Run(run func(token string)) *Token_Decode_Call[T] {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *Token_Decode_Call[T]) Return(_a0 T, _a1 error) *Token_Decode_Call[T] {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Token_Decode_Call[T]) RunAndReturn(run func(string) (T, error)) *Token_Decode_Call[T] {
	_c.Call.Return(run)
	return _c
}

// Generate provides a mock function with given fields: data
func (_m *Token[T]) Generate(data T) (string, error) {
	ret := _m.Called(data)

	if len(ret) == 0 {
		panic("no return value specified for Generate")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(T) (string, error)); ok {
		return rf(data)
	}
	if rf, ok := ret.Get(0).(func(T) string); ok {
		r0 = rf(data)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(T) error); ok {
		r1 = rf(data)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Token_Generate_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Generate'
type Token_Generate_Call[T interface{}] struct {
	*mock.Call
}

// Generate is a helper method to define mock.On call
//   - data T
func (_e *Token_Expecter[T]) Generate(data interface{}) *Token_Generate_Call[T] {
	return &Token_Generate_Call[T]{Call: _e.mock.On("Generate", data)}
}

func (_c *Token_Generate_Call[T]) Run(run func(data T)) *Token_Generate_Call[T] {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(T))
	})
	return _c
}

func (_c *Token_Generate_Call[T]) Return(_a0 string, _a1 error) *Token_Generate_Call[T] {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Token_Generate_Call[T]) RunAndReturn(run func(T) (string, error)) *Token_Generate_Call[T] {
	_c.Call.Return(run)
	return _c
}

// NewToken creates a new instance of Token. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewToken[T interface{}](t interface {
	mock.TestingT
	Cleanup(func())
}) *Token[T] {
	mock := &Token[T]{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
