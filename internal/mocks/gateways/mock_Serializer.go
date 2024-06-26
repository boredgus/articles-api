// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// Serializer is an autogenerated mock type for the Serializer type
type Serializer[T interface{}] struct {
	mock.Mock
}

type Serializer_Expecter[T interface{}] struct {
	mock *mock.Mock
}

func (_m *Serializer[T]) EXPECT() *Serializer_Expecter[T] {
	return &Serializer_Expecter[T]{mock: &_m.Mock}
}

// Deserialize provides a mock function with given fields: value
func (_m *Serializer[T]) Deserialize(value string) ([]T, error) {
	ret := _m.Called(value)

	if len(ret) == 0 {
		panic("no return value specified for Deserialize")
	}

	var r0 []T
	var r1 error
	if rf, ok := ret.Get(0).(func(string) ([]T, error)); ok {
		return rf(value)
	}
	if rf, ok := ret.Get(0).(func(string) []T); ok {
		r0 = rf(value)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]T)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(value)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Serializer_Deserialize_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Deserialize'
type Serializer_Deserialize_Call[T interface{}] struct {
	*mock.Call
}

// Deserialize is a helper method to define mock.On call
//   - value string
func (_e *Serializer_Expecter[T]) Deserialize(value interface{}) *Serializer_Deserialize_Call[T] {
	return &Serializer_Deserialize_Call[T]{Call: _e.mock.On("Deserialize", value)}
}

func (_c *Serializer_Deserialize_Call[T]) Run(run func(value string)) *Serializer_Deserialize_Call[T] {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *Serializer_Deserialize_Call[T]) Return(_a0 []T, _a1 error) *Serializer_Deserialize_Call[T] {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Serializer_Deserialize_Call[T]) RunAndReturn(run func(string) ([]T, error)) *Serializer_Deserialize_Call[T] {
	_c.Call.Return(run)
	return _c
}

// Serialize provides a mock function with given fields: value
func (_m *Serializer[T]) Serialize(value []T) (string, error) {
	ret := _m.Called(value)

	if len(ret) == 0 {
		panic("no return value specified for Serialize")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func([]T) (string, error)); ok {
		return rf(value)
	}
	if rf, ok := ret.Get(0).(func([]T) string); ok {
		r0 = rf(value)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func([]T) error); ok {
		r1 = rf(value)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Serializer_Serialize_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Serialize'
type Serializer_Serialize_Call[T interface{}] struct {
	*mock.Call
}

// Serialize is a helper method to define mock.On call
//   - value []T
func (_e *Serializer_Expecter[T]) Serialize(value interface{}) *Serializer_Serialize_Call[T] {
	return &Serializer_Serialize_Call[T]{Call: _e.mock.On("Serialize", value)}
}

func (_c *Serializer_Serialize_Call[T]) Run(run func(value []T)) *Serializer_Serialize_Call[T] {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].([]T))
	})
	return _c
}

func (_c *Serializer_Serialize_Call[T]) Return(_a0 string, _a1 error) *Serializer_Serialize_Call[T] {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Serializer_Serialize_Call[T]) RunAndReturn(run func([]T) (string, error)) *Serializer_Serialize_Call[T] {
	_c.Call.Return(run)
	return _c
}

// NewSerializer creates a new instance of Serializer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewSerializer[T interface{}](t interface {
	mock.TestingT
	Cleanup(func())
}) *Serializer[T] {
	mock := &Serializer[T]{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
