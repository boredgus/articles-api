// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	domain "a-article/internal/domain"

	mock "github.com/stretchr/testify/mock"

	repo "a-article/internal/models/repo"
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
func (_m *UserRepository) Create(user repo.User) error {
	ret := _m.Called(user)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(repo.User) error); ok {
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
//   - user repo.User
func (_e *UserRepository_Expecter) Create(user interface{}) *UserRepository_Create_Call {
	return &UserRepository_Create_Call{Call: _e.mock.On("Create", user)}
}

func (_c *UserRepository_Create_Call) Run(run func(user repo.User)) *UserRepository_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(repo.User))
	})
	return _c
}

func (_c *UserRepository_Create_Call) Return(_a0 error) *UserRepository_Create_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *UserRepository_Create_Call) RunAndReturn(run func(repo.User) error) *UserRepository_Create_Call {
	_c.Call.Return(run)
	return _c
}

// Delete provides a mock function with given fields: oid
func (_m *UserRepository) Delete(oid string) error {
	ret := _m.Called(oid)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(oid)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UserRepository_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type UserRepository_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - oid string
func (_e *UserRepository_Expecter) Delete(oid interface{}) *UserRepository_Delete_Call {
	return &UserRepository_Delete_Call{Call: _e.mock.On("Delete", oid)}
}

func (_c *UserRepository_Delete_Call) Run(run func(oid string)) *UserRepository_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *UserRepository_Delete_Call) Return(_a0 error) *UserRepository_Delete_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *UserRepository_Delete_Call) RunAndReturn(run func(string) error) *UserRepository_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// Get provides a mock function with given fields: username
func (_m *UserRepository) Get(username string) (repo.User, error) {
	ret := _m.Called(username)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 repo.User
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (repo.User, error)); ok {
		return rf(username)
	}
	if rf, ok := ret.Get(0).(func(string) repo.User); ok {
		r0 = rf(username)
	} else {
		r0 = ret.Get(0).(repo.User)
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

func (_c *UserRepository_Get_Call) Return(_a0 repo.User, _a1 error) *UserRepository_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UserRepository_Get_Call) RunAndReturn(run func(string) (repo.User, error)) *UserRepository_Get_Call {
	_c.Call.Return(run)
	return _c
}

// GetByOId provides a mock function with given fields: oid
func (_m *UserRepository) GetByOId(oid string) (repo.User, error) {
	ret := _m.Called(oid)

	if len(ret) == 0 {
		panic("no return value specified for GetByOId")
	}

	var r0 repo.User
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (repo.User, error)); ok {
		return rf(oid)
	}
	if rf, ok := ret.Get(0).(func(string) repo.User); ok {
		r0 = rf(oid)
	} else {
		r0 = ret.Get(0).(repo.User)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(oid)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserRepository_GetByOId_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetByOId'
type UserRepository_GetByOId_Call struct {
	*mock.Call
}

// GetByOId is a helper method to define mock.On call
//   - oid string
func (_e *UserRepository_Expecter) GetByOId(oid interface{}) *UserRepository_GetByOId_Call {
	return &UserRepository_GetByOId_Call{Call: _e.mock.On("GetByOId", oid)}
}

func (_c *UserRepository_GetByOId_Call) Run(run func(oid string)) *UserRepository_GetByOId_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *UserRepository_GetByOId_Call) Return(_a0 repo.User, _a1 error) *UserRepository_GetByOId_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UserRepository_GetByOId_Call) RunAndReturn(run func(string) (repo.User, error)) *UserRepository_GetByOId_Call {
	_c.Call.Return(run)
	return _c
}

// GetSignupRequest provides a mock function with given fields: email
func (_m *UserRepository) GetSignupRequest(email string) (repo.SignupRequest, error) {
	ret := _m.Called(email)

	if len(ret) == 0 {
		panic("no return value specified for GetSignupRequest")
	}

	var r0 repo.SignupRequest
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (repo.SignupRequest, error)); ok {
		return rf(email)
	}
	if rf, ok := ret.Get(0).(func(string) repo.SignupRequest); ok {
		r0 = rf(email)
	} else {
		r0 = ret.Get(0).(repo.SignupRequest)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserRepository_GetSignupRequest_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetSignupRequest'
type UserRepository_GetSignupRequest_Call struct {
	*mock.Call
}

// GetSignupRequest is a helper method to define mock.On call
//   - email string
func (_e *UserRepository_Expecter) GetSignupRequest(email interface{}) *UserRepository_GetSignupRequest_Call {
	return &UserRepository_GetSignupRequest_Call{Call: _e.mock.On("GetSignupRequest", email)}
}

func (_c *UserRepository_GetSignupRequest_Call) Run(run func(email string)) *UserRepository_GetSignupRequest_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *UserRepository_GetSignupRequest_Call) Return(_a0 repo.SignupRequest, _a1 error) *UserRepository_GetSignupRequest_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UserRepository_GetSignupRequest_Call) RunAndReturn(run func(string) (repo.SignupRequest, error)) *UserRepository_GetSignupRequest_Call {
	_c.Call.Return(run)
	return _c
}

// RegisterSignupRequest provides a mock function with given fields: request
func (_m *UserRepository) RegisterSignupRequest(request repo.SignupRequest) error {
	ret := _m.Called(request)

	if len(ret) == 0 {
		panic("no return value specified for RegisterSignupRequest")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(repo.SignupRequest) error); ok {
		r0 = rf(request)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UserRepository_RegisterSignupRequest_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RegisterSignupRequest'
type UserRepository_RegisterSignupRequest_Call struct {
	*mock.Call
}

// RegisterSignupRequest is a helper method to define mock.On call
//   - request repo.SignupRequest
func (_e *UserRepository_Expecter) RegisterSignupRequest(request interface{}) *UserRepository_RegisterSignupRequest_Call {
	return &UserRepository_RegisterSignupRequest_Call{Call: _e.mock.On("RegisterSignupRequest", request)}
}

func (_c *UserRepository_RegisterSignupRequest_Call) Run(run func(request repo.SignupRequest)) *UserRepository_RegisterSignupRequest_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(repo.SignupRequest))
	})
	return _c
}

func (_c *UserRepository_RegisterSignupRequest_Call) Return(_a0 error) *UserRepository_RegisterSignupRequest_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *UserRepository_RegisterSignupRequest_Call) RunAndReturn(run func(repo.SignupRequest) error) *UserRepository_RegisterSignupRequest_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateRole provides a mock function with given fields: oid, role
func (_m *UserRepository) UpdateRole(oid string, role domain.UserRole) error {
	ret := _m.Called(oid, role)

	if len(ret) == 0 {
		panic("no return value specified for UpdateRole")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, domain.UserRole) error); ok {
		r0 = rf(oid, role)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UserRepository_UpdateRole_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateRole'
type UserRepository_UpdateRole_Call struct {
	*mock.Call
}

// UpdateRole is a helper method to define mock.On call
//   - oid string
//   - role domain.UserRole
func (_e *UserRepository_Expecter) UpdateRole(oid interface{}, role interface{}) *UserRepository_UpdateRole_Call {
	return &UserRepository_UpdateRole_Call{Call: _e.mock.On("UpdateRole", oid, role)}
}

func (_c *UserRepository_UpdateRole_Call) Run(run func(oid string, role domain.UserRole)) *UserRepository_UpdateRole_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(domain.UserRole))
	})
	return _c
}

func (_c *UserRepository_UpdateRole_Call) Return(_a0 error) *UserRepository_UpdateRole_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *UserRepository_UpdateRole_Call) RunAndReturn(run func(string, domain.UserRole) error) *UserRepository_UpdateRole_Call {
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
