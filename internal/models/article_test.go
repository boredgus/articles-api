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
			Create(mock.Anything, mock.Anything).Return(res.createErr).Once()
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
		username string
		article  *domain.Article
	}
	type mockedRes struct {
		isOwnerErr error
		updateErr  error
	}
	repoMock := repoMocks.NewArticleRepository(t)
	setup := func(res mockedRes) func() {
		isOwnerCall := repoMock.EXPECT().IsOwner(mock.Anything, mock.Anything).
			Return(domain.Article{}, res.isOwnerErr).Once()
		updateCall := repoMock.EXPECT().
			Update(mock.Anything, mock.Anything).NotBefore(isOwnerCall).
			Return(res.updateErr).Once()
		return func() {
			isOwnerCall.Unset()
			updateCall.Unset()
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
			name:      "user is not an owner of article",
			args:      args{article: &domain.Article{Theme: ""}},
			mockedRes: mockedRes{isOwnerErr: UserIsNotAnOwnerErr},
			wantErr:   UserIsNotAnOwnerErr,
		},
		{
			name:    "article data is not valid",
			args:    args{article: &domain.Article{Theme: ""}},
			wantErr: InvalidArticleErr,
		},
		{
			name:      "failed to update article",
			args:      args{article: &domain.Article{Theme: "t", Tags: []string{}}},
			mockedRes: mockedRes{updateErr: someError},
			wantErr:   someError,
		},
		{
			name: "success",
			args: args{article: &domain.Article{Theme: "t", Tags: []string{}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanSetup := setup(tt.mockedRes)
			defer cleanSetup()
			err := NewArticleModel(repoMock).Update(tt.args.username, tt.args.article)
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
