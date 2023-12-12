package models

import (
	"errors"
	"testing"
	"user-management/internal/domain"
	repoMocks "user-management/internal/mocks/repo"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestArticleService_Create(t *testing.T) {
	type args struct {
		userOId string
		article *domain.Article
	}
	type mockedRes struct {
		createErr error
	}
	repoMock := repoMocks.NewArticleRepository(t)
	setup := func(res mockedRes) func() {
		repoCall := repoMock.EXPECT().
			CreateArticle(mock.Anything, mock.Anything).Return(res.createErr).Once()
		return func() {
			repoCall.Unset()
		}
	}
	someError := errors.New("some error")
	tests := []struct {
		name      string
		mockedRes mockedRes
		args      args
		wantErr   error
	}{
		{
			name:    "article data is invalid",
			args:    args{article: &domain.Article{Theme: ""}},
			wantErr: InvalidArticleErr,
		},
		{
			name:      "there is no user with such oid",
			args:      args{article: &domain.Article{Theme: "error"}},
			mockedRes: mockedRes{createErr: ArticleNotFoundErr},
			wantErr:   UnknownUserErr,
		},
		{
			name:      "failed to insert article data to db",
			args:      args{article: &domain.Article{Theme: "error"}},
			mockedRes: mockedRes{createErr: someError},
			wantErr:   someError,
		},
		{
			name:    "success",
			args:    args{article: &domain.Article{Theme: "success"}},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanSetup := setup(tt.mockedRes)
			defer cleanSetup()
			err := NewArticleModel(repoMock).Create(tt.args.userOId, tt.args.article)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			assert.Nil(t, err)
		})
	}
}

func TestArticleService_GetForUser(t *testing.T) {
	type args struct {
		username string
		page     int
		limit    int
	}
	type mockedRes struct {
		getErr      error
		getArticles []domain.Article
	}
	repoMock := repoMocks.NewArticleRepository(t)
	setup := func(res mockedRes) func() {
		repoCall := repoMock.EXPECT().
			GetForUser(mock.Anything, mock.Anything, mock.Anything).
			Return(res.getArticles, res.getErr).Once()
		return func() {
			repoCall.Unset()
		}
	}
	someError := errors.New("some error")
	articles := []domain.Article{{Theme: "theme1"}, {Theme: "theme2"}, {Theme: "theme3"}}
	tests := []struct {
		name           string
		args           args
		mockedRes      mockedRes
		wantArticle    []domain.Article
		wantPagination PaginationData
		wantErr        error
	}{
		{
			name:           "failed to get articles",
			mockedRes:      mockedRes{getArticles: []domain.Article{}, getErr: someError},
			wantArticle:    nil,
			wantPagination: PaginationData{},
			wantErr:        someError,
		},
		{
			name:           "success",
			args:           args{limit: 10},
			mockedRes:      mockedRes{getArticles: articles, getErr: nil},
			wantArticle:    articles,
			wantPagination: PaginationData{Page: 0, Limit: 10, Count: len(articles)},
			wantErr:        nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanSetup := setup(tt.mockedRes)
			defer cleanSetup()
			articles, pagination, err := NewArticleModel(repoMock).GetForUser(tt.args.username, tt.args.page, tt.args.limit)
			assert.Equal(t, articles, tt.wantArticle, "articles_check")
			assert.Equal(t, pagination, tt.wantPagination, "pagination_check")
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr, "error_check")
				return
			}
			assert.Nil(t, err)
		})
	}
}

func TestArticleService_Update(t *testing.T) {
	type args struct {
		userOId  string
		userRole string
		article  *domain.Article
	}
	type mockedRes struct {
		oldArticle    domain.Article
		oldArticleErr error
		isOwnerErr    error
		updateErr     error
		addTagsErr    error
		removeTagsErr error
	}
	repoMock := repoMocks.NewArticleRepository(t)
	setup := func(res mockedRes) func() {
		getCall := repoMock.EXPECT().Get(mock.Anything).Return(res.oldArticle, res.oldArticleErr).Once()
		updateCall := repoMock.EXPECT().
			UpdateArticle(mock.Anything, mock.Anything, mock.Anything).NotBefore(getCall).
			Return(res.updateErr).Once()
		calls := []*mock.Call{
			getCall,
			repoMock.EXPECT().IsOwner(mock.Anything, mock.Anything).
				Return(res.isOwnerErr).NotBefore(getCall).Maybe(),
			updateCall,
			repoMock.EXPECT().RemoveTagsFromArticle(mock.Anything, mock.Anything).
				NotBefore(updateCall).Return(res.removeTagsErr).Maybe(),
			repoMock.EXPECT().AddTagsForArticle(mock.Anything, mock.Anything).
				NotBefore(updateCall).Return(res.addTagsErr).Maybe(),
		}
		return func() {
			for _, call := range calls {
				call.Unset()
			}
		}
	}
	someError := errors.New("some error")
	tests := []struct {
		name      string
		args      args
		mockedRes mockedRes
		wantErr   error
	}{
		{
			name:      "failed to get old article",
			args:      args{article: &domain.Article{}},
			mockedRes: mockedRes{oldArticleErr: someError},
			wantErr:   someError,
		},
		{
			name:    "not enough rights to update article",
			args:    args{userRole: "kek", article: &domain.Article{}},
			wantErr: NotEnoughRightsErr,
		},
		{
			name:    "article data is not valid",
			args:    args{userRole: "admin", article: &domain.Article{Theme: ""}},
			wantErr: InvalidArticleErr,
		},
		{
			name:      "failed to update article",
			args:      args{userRole: "admin", article: &domain.Article{Theme: "t", Tags: []string{}}},
			mockedRes: mockedRes{updateErr: someError},
			wantErr:   someError,
		},
		{
			name: "failed to remove article tags",
			args: args{userRole: "admin", article: &domain.Article{Theme: "t", Tags: []string{}}},
			mockedRes: mockedRes{
				oldArticle:    domain.Article{Theme: "t", Tags: []string{"old"}},
				removeTagsErr: someError},
			wantErr: someError,
		},
		{
			name: "failed to add article tags",
			args: args{userRole: "admin", article: &domain.Article{Theme: "t", Tags: []string{"new"}}},
			mockedRes: mockedRes{
				oldArticle: domain.Article{Theme: "t", Tags: []string{}},
				addTagsErr: someError},
			wantErr: someError,
		},
		{
			name: "success",
			args: args{userRole: "admin", article: &domain.Article{Theme: "t", Tags: []string{}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanSetup := setup(tt.mockedRes)
			defer cleanSetup()
			err := NewArticleModel(repoMock).Update(tt.args.userOId, tt.args.userRole, tt.args.article)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			assert.Nil(t, err)
		})
	}
}

func TestArticleService_Get(t *testing.T) {
	type args struct {
		articleOId string
	}
	type mockedRes struct {
		article domain.Article
		err     error
	}
	repoMock := repoMocks.NewArticleRepository(t)
	setup := func(res mockedRes) func() {
		getCall := repoMock.EXPECT().
			Get(mock.Anything).Return(res.article, res.err).Once()
		return func() {
			getCall.Unset()
		}
	}
	someErr := errors.New("some err")
	tests := []struct {
		name      string
		args      args
		mockedRes mockedRes
		wantA     domain.Article
		wantErr   error
	}{
		{
			name:      "failed to fetch article",
			mockedRes: mockedRes{article: domain.Article{}, err: someErr},
			wantA:     domain.Article{},
			wantErr:   someErr,
		},
		{
			name:      "success",
			mockedRes: mockedRes{article: domain.Article{}},
			wantA:     domain.Article{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanSetup := setup(tt.mockedRes)
			defer cleanSetup()
			got, err := NewArticleModel(repoMock).Get(tt.args.articleOId)
			assert.Equal(t, got, tt.wantA)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			assert.Nil(t, err)
		})
	}
}

func TestArticleService_Delete(t *testing.T) {
	type args struct {
		userOId    string
		userRole   string
		articleOId string
	}
	type mockedRes struct {
		getErr    error
		article   domain.Article
		deleteErr error
	}
	repoMock := repoMocks.NewArticleRepository(t)
	setup := func(res mockedRes) func() {
		getCall := repoMock.EXPECT().
			Get(mock.Anything).Return(res.article, res.getErr).Once()
		deleteCall := repoMock.EXPECT().DeleteArticle(mock.Anything, res.article.Tags).
			NotBefore(getCall).Return(res.deleteErr)
		return func() {
			getCall.Unset()
			deleteCall.Unset()
		}
	}
	someErr := errors.New("some err")
	tests := []struct {
		name      string
		args      args
		mockedRes mockedRes
		wantErr   error
	}{
		{
			name:      "failed to get article",
			args:      args{userRole: "admin"},
			mockedRes: mockedRes{getErr: someErr},
			wantErr:   someErr,
		},
		{
			name:    "not enough rights to delete article",
			args:    args{userRole: "kek"},
			wantErr: NotEnoughRightsErr,
		},
		{
			name:      "failed to delete article",
			args:      args{userRole: "admin"},
			mockedRes: mockedRes{deleteErr: someErr},
			wantErr:   someErr,
		},
		{
			name:    "success",
			args:    args{userRole: "admin"},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanSetup := setup(tt.mockedRes)
			defer cleanSetup()
			err := ArticleService{repo: repoMock}.Delete(tt.args.userOId, tt.args.userRole, tt.args.articleOId)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			assert.Nil(t, err)
		})
	}
}

func TestArticleService_checkRights(t *testing.T) {
	type args struct {
		userOId    string
		userRole   string
		articleOId string
	}
	type mockedRes struct {
		isOwnerErr error
	}
	repoMock := repoMocks.NewArticleRepository(t)
	setup := func(res mockedRes) func() {
		isOwnerCall := repoMock.EXPECT().IsOwner(mock.Anything, mock.Anything).
			Return(res.isOwnerErr).Maybe()
		return func() {
			isOwnerCall.Unset()
		}
	}
	someError := errors.New("some err")
	tests := []struct {
		name      string
		args      args
		mockedRes mockedRes
		wantErr   error
	}{
		{
			name:      "user with default role is trying to change not his article",
			args:      args{userRole: "user"},
			mockedRes: mockedRes{isOwnerErr: ArticleNotFoundErr},
			wantErr:   NotEnoughRightsErr,
		},
		{
			name:      "server error on IsOwner check",
			args:      args{userRole: "user"},
			mockedRes: mockedRes{isOwnerErr: someError},
			wantErr:   someError,
		},
		{
			name:    "server error on IsOwner check",
			args:    args{userRole: "kek"},
			wantErr: NotEnoughRightsErr,
		},
		{
			name:    "success",
			args:    args{userRole: "kek"},
			wantErr: NotEnoughRightsErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanSetup := setup(tt.mockedRes)
			defer cleanSetup()
			err := ArticleService{repo: repoMock}.checkRights(tt.args.userOId, tt.args.userRole, tt.args.articleOId)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			assert.Nil(t, err)
		})
	}
}
