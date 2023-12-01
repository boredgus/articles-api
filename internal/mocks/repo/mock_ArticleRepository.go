// Code generated by mockery v2.38.0. DO NOT EDIT.

package mocks

import (
	domain "user-management/internal/domain"

	mock "github.com/stretchr/testify/mock"

	repo "user-management/internal/models/repo"
)

// ArticleRepository is an autogenerated mock type for the ArticleRepository type
type ArticleRepository struct {
	mock.Mock
}

type ArticleRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *ArticleRepository) EXPECT() *ArticleRepository_Expecter {
	return &ArticleRepository_Expecter{mock: &_m.Mock}
}

// Create provides a mock function with given fields: userOId, article
func (_m *ArticleRepository) Create(userOId string, article repo.ArticleData) error {
	ret := _m.Called(userOId, article)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, repo.ArticleData) error); ok {
		r0 = rf(userOId, article)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ArticleRepository_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type ArticleRepository_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - userOId string
//   - article repo.ArticleData
func (_e *ArticleRepository_Expecter) Create(userOId interface{}, article interface{}) *ArticleRepository_Create_Call {
	return &ArticleRepository_Create_Call{Call: _e.mock.On("Create", userOId, article)}
}

func (_c *ArticleRepository_Create_Call) Run(run func(userOId string, article repo.ArticleData)) *ArticleRepository_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(repo.ArticleData))
	})
	return _c
}

func (_c *ArticleRepository_Create_Call) Return(_a0 error) *ArticleRepository_Create_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ArticleRepository_Create_Call) RunAndReturn(run func(string, repo.ArticleData) error) *ArticleRepository_Create_Call {
	_c.Call.Return(run)
	return _c
}

// GetForUser provides a mock function with given fields: username, page, limit
func (_m *ArticleRepository) GetForUser(username string, page int, limit int) ([]domain.Article, error) {
	ret := _m.Called(username, page, limit)

	if len(ret) == 0 {
		panic("no return value specified for GetForUser")
	}

	var r0 []domain.Article
	var r1 error
	if rf, ok := ret.Get(0).(func(string, int, int) ([]domain.Article, error)); ok {
		return rf(username, page, limit)
	}
	if rf, ok := ret.Get(0).(func(string, int, int) []domain.Article); ok {
		r0 = rf(username, page, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Article)
		}
	}

	if rf, ok := ret.Get(1).(func(string, int, int) error); ok {
		r1 = rf(username, page, limit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ArticleRepository_GetForUser_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetForUser'
type ArticleRepository_GetForUser_Call struct {
	*mock.Call
}

// GetForUser is a helper method to define mock.On call
//   - username string
//   - page int
//   - limit int
func (_e *ArticleRepository_Expecter) GetForUser(username interface{}, page interface{}, limit interface{}) *ArticleRepository_GetForUser_Call {
	return &ArticleRepository_GetForUser_Call{Call: _e.mock.On("GetForUser", username, page, limit)}
}

func (_c *ArticleRepository_GetForUser_Call) Run(run func(username string, page int, limit int)) *ArticleRepository_GetForUser_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(int), args[2].(int))
	})
	return _c
}

func (_c *ArticleRepository_GetForUser_Call) Return(_a0 []domain.Article, _a1 error) *ArticleRepository_GetForUser_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ArticleRepository_GetForUser_Call) RunAndReturn(run func(string, int, int) ([]domain.Article, error)) *ArticleRepository_GetForUser_Call {
	_c.Call.Return(run)
	return _c
}

// IsOwner provides a mock function with given fields: articleOId, username
func (_m *ArticleRepository) IsOwner(articleOId string, username string) error {
	ret := _m.Called(articleOId, username)

	if len(ret) == 0 {
		panic("no return value specified for IsOwner")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(articleOId, username)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ArticleRepository_IsOwner_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'IsOwner'
type ArticleRepository_IsOwner_Call struct {
	*mock.Call
}

// IsOwner is a helper method to define mock.On call
//   - articleOId string
//   - username string
func (_e *ArticleRepository_Expecter) IsOwner(articleOId interface{}, username interface{}) *ArticleRepository_IsOwner_Call {
	return &ArticleRepository_IsOwner_Call{Call: _e.mock.On("IsOwner", articleOId, username)}
}

func (_c *ArticleRepository_IsOwner_Call) Run(run func(articleOId string, username string)) *ArticleRepository_IsOwner_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string))
	})
	return _c
}

func (_c *ArticleRepository_IsOwner_Call) Return(_a0 error) *ArticleRepository_IsOwner_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ArticleRepository_IsOwner_Call) RunAndReturn(run func(string, string) error) *ArticleRepository_IsOwner_Call {
	_c.Call.Return(run)
	return _c
}

// Update provides a mock function with given fields: article
func (_m *ArticleRepository) Update(article repo.ArticleData) error {
	ret := _m.Called(article)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(repo.ArticleData) error); ok {
		r0 = rf(article)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ArticleRepository_Update_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Update'
type ArticleRepository_Update_Call struct {
	*mock.Call
}

// Update is a helper method to define mock.On call
//   - article repo.ArticleData
func (_e *ArticleRepository_Expecter) Update(article interface{}) *ArticleRepository_Update_Call {
	return &ArticleRepository_Update_Call{Call: _e.mock.On("Update", article)}
}

func (_c *ArticleRepository_Update_Call) Run(run func(article repo.ArticleData)) *ArticleRepository_Update_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(repo.ArticleData))
	})
	return _c
}

func (_c *ArticleRepository_Update_Call) Return(_a0 error) *ArticleRepository_Update_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ArticleRepository_Update_Call) RunAndReturn(run func(repo.ArticleData) error) *ArticleRepository_Update_Call {
	_c.Call.Return(run)
	return _c
}

// NewArticleRepository creates a new instance of ArticleRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewArticleRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *ArticleRepository {
	mock := &ArticleRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
