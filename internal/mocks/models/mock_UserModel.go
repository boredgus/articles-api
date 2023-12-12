// Code generated by mockery v2.38.0. DO NOT EDIT.

package mocks

import (
	domain "user-management/internal/domain"

	mock "github.com/stretchr/testify/mock"
)

// UserModel is an autogenerated mock type for the UserModel type
type UserModel struct {
	mock.Mock
}

type UserModel_Expecter struct {
	mock *mock.Mock
}

func (_m *UserModel) EXPECT() *UserModel_Expecter {
	return &UserModel_Expecter{mock: &_m.Mock}
}

// Authorize provides a mock function with given fields: username, password
func (_m *UserModel) Authorize(username string, password string) (string, error) {
	ret := _m.Called(username, password)

	if len(ret) == 0 {
		panic("no return value specified for Authorize")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string) (string, error)); ok {
		return rf(username, password)
	}
	if rf, ok := ret.Get(0).(func(string, string) string); ok {
		r0 = rf(username, password)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(username, password)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserModel_Authorize_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Authorize'
type UserModel_Authorize_Call struct {
	*mock.Call
}

// Authorize is a helper method to define mock.On call
//   - username string
//   - password string
func (_e *UserModel_Expecter) Authorize(username interface{}, password interface{}) *UserModel_Authorize_Call {
	return &UserModel_Authorize_Call{Call: _e.mock.On("Authorize", username, password)}
}

func (_c *UserModel_Authorize_Call) Run(run func(username string, password string)) *UserModel_Authorize_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string))
	})
	return _c
}

func (_c *UserModel_Authorize_Call) Return(_a0 string, _a1 error) *UserModel_Authorize_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UserModel_Authorize_Call) RunAndReturn(run func(string, string) (string, error)) *UserModel_Authorize_Call {
	_c.Call.Return(run)
	return _c
}

// Create provides a mock function with given fields: user
func (_m *UserModel) Create(user domain.User) error {
	ret := _m.Called(user)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(domain.User) error); ok {
		r0 = rf(user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UserModel_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type UserModel_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - user domain.User
func (_e *UserModel_Expecter) Create(user interface{}) *UserModel_Create_Call {
	return &UserModel_Create_Call{Call: _e.mock.On("Create", user)}
}

func (_c *UserModel_Create_Call) Run(run func(user domain.User)) *UserModel_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(domain.User))
	})
	return _c
}

func (_c *UserModel_Create_Call) Return(_a0 error) *UserModel_Create_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *UserModel_Create_Call) RunAndReturn(run func(domain.User) error) *UserModel_Create_Call {
	_c.Call.Return(run)
	return _c
}

// Exists provides a mock function with given fields: oid, password
func (_m *UserModel) Exists(oid string, password string) error {
	ret := _m.Called(oid, password)

	if len(ret) == 0 {
		panic("no return value specified for Exists")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(oid, password)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UserModel_Exists_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Exists'
type UserModel_Exists_Call struct {
	*mock.Call
}

// Exists is a helper method to define mock.On call
//   - oid string
//   - password string
func (_e *UserModel_Expecter) Exists(oid interface{}, password interface{}) *UserModel_Exists_Call {
	return &UserModel_Exists_Call{Call: _e.mock.On("Exists", oid, password)}
}

func (_c *UserModel_Exists_Call) Run(run func(oid string, password string)) *UserModel_Exists_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string))
	})
	return _c
}

func (_c *UserModel_Exists_Call) Return(_a0 error) *UserModel_Exists_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *UserModel_Exists_Call) RunAndReturn(run func(string, string) error) *UserModel_Exists_Call {
	_c.Call.Return(run)
	return _c
}

// NewUserModel creates a new instance of UserModel. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUserModel(t interface {
	mock.TestingT
	Cleanup(func())
}) *UserModel {
	mock := &UserModel{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
