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

// AddTagsForArticle provides a mock function with given fields: articleOId, tags
func (_m *ArticleRepository) AddTagsForArticle(articleOId string, tags []string) error {
	ret := _m.Called(articleOId, tags)

	if len(ret) == 0 {
		panic("no return value specified for AddTagsForArticle")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, []string) error); ok {
		r0 = rf(articleOId, tags)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ArticleRepository_AddTagsForArticle_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddTagsForArticle'
type ArticleRepository_AddTagsForArticle_Call struct {
	*mock.Call
}

// AddTagsForArticle is a helper method to define mock.On call
//   - articleOId string
//   - tags []string
func (_e *ArticleRepository_Expecter) AddTagsForArticle(articleOId interface{}, tags interface{}) *ArticleRepository_AddTagsForArticle_Call {
	return &ArticleRepository_AddTagsForArticle_Call{Call: _e.mock.On("AddTagsForArticle", articleOId, tags)}
}

func (_c *ArticleRepository_AddTagsForArticle_Call) Run(run func(articleOId string, tags []string)) *ArticleRepository_AddTagsForArticle_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].([]string))
	})
	return _c
}

func (_c *ArticleRepository_AddTagsForArticle_Call) Return(_a0 error) *ArticleRepository_AddTagsForArticle_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ArticleRepository_AddTagsForArticle_Call) RunAndReturn(run func(string, []string) error) *ArticleRepository_AddTagsForArticle_Call {
	_c.Call.Return(run)
	return _c
}

// CreateArticle provides a mock function with given fields: userOId, article
func (_m *ArticleRepository) CreateArticle(userOId string, article repo.ArticleData) error {
	ret := _m.Called(userOId, article)

	if len(ret) == 0 {
		panic("no return value specified for CreateArticle")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, repo.ArticleData) error); ok {
		r0 = rf(userOId, article)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ArticleRepository_CreateArticle_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateArticle'
type ArticleRepository_CreateArticle_Call struct {
	*mock.Call
}

// CreateArticle is a helper method to define mock.On call
//   - userOId string
//   - article repo.ArticleData
func (_e *ArticleRepository_Expecter) CreateArticle(userOId interface{}, article interface{}) *ArticleRepository_CreateArticle_Call {
	return &ArticleRepository_CreateArticle_Call{Call: _e.mock.On("CreateArticle", userOId, article)}
}

func (_c *ArticleRepository_CreateArticle_Call) Run(run func(userOId string, article repo.ArticleData)) *ArticleRepository_CreateArticle_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(repo.ArticleData))
	})
	return _c
}

func (_c *ArticleRepository_CreateArticle_Call) Return(_a0 error) *ArticleRepository_CreateArticle_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ArticleRepository_CreateArticle_Call) RunAndReturn(run func(string, repo.ArticleData) error) *ArticleRepository_CreateArticle_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteArticle provides a mock function with given fields: oid, tags
func (_m *ArticleRepository) DeleteArticle(oid string, tags []string) error {
	ret := _m.Called(oid, tags)

	if len(ret) == 0 {
		panic("no return value specified for DeleteArticle")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, []string) error); ok {
		r0 = rf(oid, tags)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ArticleRepository_DeleteArticle_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteArticle'
type ArticleRepository_DeleteArticle_Call struct {
	*mock.Call
}

// DeleteArticle is a helper method to define mock.On call
//   - oid string
//   - tags []string
func (_e *ArticleRepository_Expecter) DeleteArticle(oid interface{}, tags interface{}) *ArticleRepository_DeleteArticle_Call {
	return &ArticleRepository_DeleteArticle_Call{Call: _e.mock.On("DeleteArticle", oid, tags)}
}

func (_c *ArticleRepository_DeleteArticle_Call) Run(run func(oid string, tags []string)) *ArticleRepository_DeleteArticle_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].([]string))
	})
	return _c
}

func (_c *ArticleRepository_DeleteArticle_Call) Return(_a0 error) *ArticleRepository_DeleteArticle_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ArticleRepository_DeleteArticle_Call) RunAndReturn(run func(string, []string) error) *ArticleRepository_DeleteArticle_Call {
	_c.Call.Return(run)
	return _c
}

// Get provides a mock function with given fields: articleOId
func (_m *ArticleRepository) Get(articleOId string) (domain.Article, error) {
	ret := _m.Called(articleOId)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 domain.Article
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (domain.Article, error)); ok {
		return rf(articleOId)
	}
	if rf, ok := ret.Get(0).(func(string) domain.Article); ok {
		r0 = rf(articleOId)
	} else {
		r0 = ret.Get(0).(domain.Article)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(articleOId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ArticleRepository_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type ArticleRepository_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - articleOId string
func (_e *ArticleRepository_Expecter) Get(articleOId interface{}) *ArticleRepository_Get_Call {
	return &ArticleRepository_Get_Call{Call: _e.mock.On("Get", articleOId)}
}

func (_c *ArticleRepository_Get_Call) Run(run func(articleOId string)) *ArticleRepository_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *ArticleRepository_Get_Call) Return(_a0 domain.Article, _a1 error) *ArticleRepository_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ArticleRepository_Get_Call) RunAndReturn(run func(string) (domain.Article, error)) *ArticleRepository_Get_Call {
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

// IsOwner provides a mock function with given fields: articleOId, userOId
func (_m *ArticleRepository) IsOwner(articleOId string, userOId string) error {
	ret := _m.Called(articleOId, userOId)

	if len(ret) == 0 {
		panic("no return value specified for IsOwner")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(articleOId, userOId)
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
//   - userOId string
func (_e *ArticleRepository_Expecter) IsOwner(articleOId interface{}, userOId interface{}) *ArticleRepository_IsOwner_Call {
	return &ArticleRepository_IsOwner_Call{Call: _e.mock.On("IsOwner", articleOId, userOId)}
}

func (_c *ArticleRepository_IsOwner_Call) Run(run func(articleOId string, userOId string)) *ArticleRepository_IsOwner_Call {
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

// RemoveTagsFromArticle provides a mock function with given fields: articleOId, tags
func (_m *ArticleRepository) RemoveTagsFromArticle(articleOId string, tags []string) error {
	ret := _m.Called(articleOId, tags)

	if len(ret) == 0 {
		panic("no return value specified for RemoveTagsFromArticle")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, []string) error); ok {
		r0 = rf(articleOId, tags)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ArticleRepository_RemoveTagsFromArticle_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RemoveTagsFromArticle'
type ArticleRepository_RemoveTagsFromArticle_Call struct {
	*mock.Call
}

// RemoveTagsFromArticle is a helper method to define mock.On call
//   - articleOId string
//   - tags []string
func (_e *ArticleRepository_Expecter) RemoveTagsFromArticle(articleOId interface{}, tags interface{}) *ArticleRepository_RemoveTagsFromArticle_Call {
	return &ArticleRepository_RemoveTagsFromArticle_Call{Call: _e.mock.On("RemoveTagsFromArticle", articleOId, tags)}
}

func (_c *ArticleRepository_RemoveTagsFromArticle_Call) Run(run func(articleOId string, tags []string)) *ArticleRepository_RemoveTagsFromArticle_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].([]string))
	})
	return _c
}

func (_c *ArticleRepository_RemoveTagsFromArticle_Call) Return(_a0 error) *ArticleRepository_RemoveTagsFromArticle_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ArticleRepository_RemoveTagsFromArticle_Call) RunAndReturn(run func(string, []string) error) *ArticleRepository_RemoveTagsFromArticle_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateArticle provides a mock function with given fields: oid, theme, text
func (_m *ArticleRepository) UpdateArticle(oid string, theme string, text string) error {
	ret := _m.Called(oid, theme, text)

	if len(ret) == 0 {
		panic("no return value specified for UpdateArticle")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, string) error); ok {
		r0 = rf(oid, theme, text)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ArticleRepository_UpdateArticle_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateArticle'
type ArticleRepository_UpdateArticle_Call struct {
	*mock.Call
}

// UpdateArticle is a helper method to define mock.On call
//   - oid string
//   - theme string
//   - text string
func (_e *ArticleRepository_Expecter) UpdateArticle(oid interface{}, theme interface{}, text interface{}) *ArticleRepository_UpdateArticle_Call {
	return &ArticleRepository_UpdateArticle_Call{Call: _e.mock.On("UpdateArticle", oid, theme, text)}
}

func (_c *ArticleRepository_UpdateArticle_Call) Run(run func(oid string, theme string, text string)) *ArticleRepository_UpdateArticle_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *ArticleRepository_UpdateArticle_Call) Return(_a0 error) *ArticleRepository_UpdateArticle_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ArticleRepository_UpdateArticle_Call) RunAndReturn(run func(string, string, string) error) *ArticleRepository_UpdateArticle_Call {
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
