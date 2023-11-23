// Code generated by mockery v2.38.0. DO NOT EDIT.

package mocks

import (
	models "user-management/internal/models"

	mock "github.com/stretchr/testify/mock"
)

// UserRepository is an autogenerated mock type for the UserRepository type
type UserRepository struct {
	mock.Mock
}

type UserRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *UserRepository) EXPECT() *UserRepository_Expecter {
	return &UserRepository_Expecter{mock: &_m.Mock}
}

// Create provides a mock function with given fields: user
func (_m *UserRepository) Create(user models.User) error {
	ret := _m.Called(user)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(models.User) error); ok {
		r0 = rf(user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UserRepository_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type UserRepository_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - user models.User
func (_e *UserRepository_Expecter) Create(user interface{}) *UserRepository_Create_Call {
	return &UserRepository_Create_Call{Call: _e.mock.On("Create", user)}
}

func (_c *UserRepository_Create_Call) Run(run func(user models.User)) *UserRepository_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(models.User))
	})
	return _c
}

func (_c *UserRepository_Create_Call) Return(_a0 error) *UserRepository_Create_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *UserRepository_Create_Call) RunAndReturn(run func(models.User) error) *UserRepository_Create_Call {
	_c.Call.Return(run)
	return _c
}

// Get provides a mock function with given fields: username
func (_m *UserRepository) Get(username string) (models.User, error) {
	ret := _m.Called(username)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 models.User
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (models.User, error)); ok {
		return rf(username)
	}
	if rf, ok := ret.Get(0).(func(string) models.User); ok {
		r0 = rf(username)
	} else {
		r0 = ret.Get(0).(models.User)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(username)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserRepository_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type UserRepository_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - username string
func (_e *UserRepository_Expecter) Get(username interface{}) *UserRepository_Get_Call {
	return &UserRepository_Get_Call{Call: _e.mock.On("Get", username)}
}

func (_c *UserRepository_Get_Call) Run(run func(username string)) *UserRepository_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *UserRepository_Get_Call) Return(_a0 models.User, _a1 error) *UserRepository_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UserRepository_Get_Call) RunAndReturn(run func(string) (models.User, error)) *UserRepository_Get_Call {
	_c.Call.Return(run)
	return _c
}

// NewUserRepository creates a new instance of UserRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUserRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *UserRepository {
	mock := &UserRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
