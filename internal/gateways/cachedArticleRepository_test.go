package gateways

import (
	"fmt"
	"testing"
	"user-management/internal/domain"
	gtwMocks "user-management/internal/mocks/gateways"
	repoMocks "user-management/internal/mocks/repo"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCachedArticleRepository_GetArticle(t *testing.T) {
	type args struct {
		articleOId string
	}
	type mockedRes struct {
		getCacheErr   error
		setCacheErr   error
		article       domain.Article
		serializerErr error
	}
	repoMock := repoMocks.NewArticleRepository(t)
	cacheStoreMock := gtwMocks.NewCacheStore(t)
	serializerMock := gtwMocks.NewSerializer[domain.Article](t)
	setup := func(res mockedRes) func() {
		getCacheCall := cacheStoreMock.EXPECT().
			Get(mock.Anything, mock.Anything).Return(res.getCacheErr).Maybe()
		repoCall := repoMock.EXPECT().
			GetArticle(mock.Anything).NotBefore(getCacheCall).
			Return(res.article, nil).Maybe()
		serializeCall := serializerMock.EXPECT().
			Serialize(mock.Anything).NotBefore(repoCall).Return("", res.serializerErr)
		setCacheCall := cacheStoreMock.EXPECT().
			Set(mock.Anything, mock.Anything).NotBefore(serializeCall).Maybe().Return(res.setCacheErr)
		deserializeCall := serializerMock.EXPECT().
			Deserialize(mock.Anything).NotBefore(getCacheCall).Return([]domain.Article{{}}, res.serializerErr)
		return func() {
			getCacheCall.Unset()
			repoCall.Unset()
			setCacheCall.Unset()
			serializeCall.Unset()
			deserializeCall.Unset()
		}
	}
	article := domain.Article{}
	someErr := fmt.Errorf("some err")
	tests := []struct {
		name      string
		args      args
		mockedRes mockedRes
		want      domain.Article
		wantErr   error
	}{
		{
			name:      "data fetched from cache",
			mockedRes: mockedRes{article: article},
		},
		{
			name:      "failed to fetch data from cache",
			mockedRes: mockedRes{article: article, serializerErr: someErr},
			wantErr:   someErr,
		},
		{
			name: "data fetched from db and failed to serialize for cache",
			mockedRes: mockedRes{
				getCacheErr:   fmt.Errorf(""),
				article:       article,
				serializerErr: fmt.Errorf(""),
			},
			want: article,
		},
		{
			name: "data fetched from db and set to cache",
			mockedRes: mockedRes{
				getCacheErr: fmt.Errorf(""),
				article:     article,
			},
			want: article,
		},
		{
			name: "data fetched from db and failed set to cache",
			mockedRes: mockedRes{
				getCacheErr: fmt.Errorf(""),
				article:     article,
				setCacheErr: fmt.Errorf(""),
			},
			want: article,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanSetup := setup(tt.mockedRes)
			defer cleanSetup()

			got, err := (&CachedArticleRepository{
				ArticleRepository: repoMock,
				cache:             cacheStoreMock,
				serializer:        serializerMock,
			}).GetArticle(tt.args.articleOId)
			assert.Equal(t, got, tt.want)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			assert.Nil(t, err)
		})
	}
}

func TestCachedArticleRepository_GetForUser(t *testing.T) {
	type args struct {
		username string
		page     int
		limit    int
	}
	type mockedRes struct {
		getCacheErr   error
		articles      []domain.Article
		setCacheErr   error
		serializerErr error
	}
	repoMock := repoMocks.NewArticleRepository(t)
	cacheStoreMock := gtwMocks.NewCacheStore(t)
	serializerMock := gtwMocks.NewSerializer[domain.Article](t)
	articles := []domain.Article{{}}
	setup := func(res mockedRes) func() {
		getCacheCall := cacheStoreMock.EXPECT().
			Get(mock.Anything, mock.Anything).Return(res.getCacheErr).Maybe()
		repoCall := repoMock.EXPECT().
			GetForUser(mock.Anything, mock.Anything, mock.Anything).NotBefore(getCacheCall).
			Return(res.articles, nil).Maybe()
		serializeCall := serializerMock.EXPECT().
			Serialize(mock.Anything).NotBefore(repoCall).Return("", res.serializerErr)
		setCacheCall := cacheStoreMock.EXPECT().
			Set(mock.Anything, mock.Anything).NotBefore(serializeCall).Maybe().Return(res.setCacheErr)
		deserializeCall := serializerMock.EXPECT().
			Deserialize(mock.Anything).NotBefore(getCacheCall).Return(articles, res.serializerErr)
		return func() {
			getCacheCall.Unset()
			repoCall.Unset()
			setCacheCall.Unset()
			serializeCall.Unset()
			deserializeCall.Unset()
		}
	}
	tests := []struct {
		name      string
		args      args
		mockedRes mockedRes
		want      []domain.Article
		wantErr   error
	}{
		{
			name:      "data fetched from cache",
			mockedRes: mockedRes{},
			want:      articles,
		},
		{
			name: "data fetched from db and failed to serialize for cache",
			mockedRes: mockedRes{
				getCacheErr:   fmt.Errorf(""),
				articles:      articles,
				serializerErr: fmt.Errorf(""),
			},
			want: articles,
		},
		{
			name: "data fetched from db and set to cache",
			mockedRes: mockedRes{
				getCacheErr: fmt.Errorf(""),
				articles:    articles,
			},
			want: articles,
		},
		{
			name: "data fetched from db and failed set to cache",
			mockedRes: mockedRes{
				getCacheErr: fmt.Errorf(""),
				articles:    articles,
				setCacheErr: fmt.Errorf(""),
			},
			want: articles,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanSetup := setup(tt.mockedRes)
			defer cleanSetup()

			got, err := (&CachedArticleRepository{
				ArticleRepository: repoMock,
				cache:             cacheStoreMock,
				serializer:        serializerMock,
			}).GetForUser(tt.args.username, tt.args.page, tt.args.limit)
			assert.Equal(t, got, tt.want)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			assert.Nil(t, err)
		})
	}
}
