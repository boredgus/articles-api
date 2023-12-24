package models

import (
	"errors"
	"testing"
	"user-management/internal/domain"
	repoMocks "user-management/internal/mocks/repo"
	"user-management/internal/models/repo"

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
			wantErr: InvalidDataErr,
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
		articles     []domain.Article
		articlesErr  error
		reactions    repo.ArticleReactions
		reactionsErr error
	}
	repoMock := repoMocks.NewArticleRepository(t)
	setup := func(res mockedRes) func() {
		getCall := repoMock.EXPECT().
			GetForUser(mock.Anything, mock.Anything, mock.Anything).
			Return(res.articles, res.articlesErr).Once()
		reactionsCall := repoMock.EXPECT().
			GetReactionsFor("").NotBefore(getCall).Return(res.reactions, res.reactionsErr).Once()

		return func() {
			getCall.Unset()
			reactionsCall.Unset()
		}
	}
	someError := errors.New("some error")
	articles := []domain.Article{{Theme: "theme1"}}
	reactions := repo.ArticleReactions{"": {"🇫🇮": 1}}
	articlesWithReactions := []domain.Article{{Theme: "theme1", Reactions: reactions[""]}}
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
			mockedRes:      mockedRes{articles: []domain.Article{}, articlesErr: someError},
			wantArticle:    nil,
			wantPagination: PaginationData{},
			wantErr:        someError,
		},
		{
			name:           "no articles for this user",
			args:           args{limit: 10},
			mockedRes:      mockedRes{articles: []domain.Article{}},
			wantArticle:    []domain.Article{},
			wantPagination: PaginationData{Page: 0, Limit: 10, Count: 0},
			wantErr:        nil,
		},
		{
			name: "failed to get reactions",
			args: args{limit: 10},
			mockedRes: mockedRes{
				articles:     articles,
				reactionsErr: someError},
			wantPagination: PaginationData{},
			wantErr:        someError,
		},
		{
			name: "at least one article exists",
			args: args{limit: 10},
			mockedRes: mockedRes{
				articles:  articles,
				reactions: reactions},
			wantArticle:    articlesWithReactions,
			wantPagination: PaginationData{Page: 0, Limit: 10, Count: len(articlesWithReactions)},
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
		getCall := repoMock.EXPECT().GetArticle(mock.Anything).Return(res.oldArticle, res.oldArticleErr).Once()
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
			wantErr: InvalidDataErr,
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
		article      domain.Article
		getErr       error
		reactions    repo.ArticleReactions
		reactionsErr error
	}
	repoMock := repoMocks.NewArticleRepository(t)
	setup := func(res mockedRes) func() {
		getCall := repoMock.EXPECT().
			GetArticle(mock.Anything).Return(res.article, res.getErr).Once()
		reactionsCall := repoMock.EXPECT().
			GetReactionsFor("").NotBefore(getCall).Return(res.reactions, res.reactionsErr).Once()
		return func() {
			getCall.Unset()
			reactionsCall.Unset()
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
			name: "failed to fetch article",
			mockedRes: mockedRes{
				article: domain.Article{},
				getErr:  someErr},
			wantA:   domain.Article{},
			wantErr: someErr,
		},
		{
			name: "failed to fetch reactions",
			mockedRes: mockedRes{
				article:      domain.Article{},
				reactionsErr: someErr},
			wantA:   domain.Article{},
			wantErr: someErr,
		},
		{
			name: "success",
			mockedRes: mockedRes{
				article:   domain.Article{},
				reactions: repo.ArticleReactions{"": domain.ArticleReactions{"🇫🇮": 1}}},
			wantA: domain.Article{Reactions: domain.ArticleReactions{"🇫🇮": 1}},
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
			GetArticle(mock.Anything).Return(res.article, res.getErr).Once()
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
			err := NewArticleModel(repoMock).Delete(tt.args.userOId, tt.args.userRole, tt.args.articleOId)
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
			args:      args{userRole: string(domain.DefaultUserRole)},
			mockedRes: mockedRes{isOwnerErr: NotFoundErr},
			wantErr:   NotEnoughRightsErr,
		},
		{
			name:      "server error on IsOwner check",
			args:      args{userRole: string(domain.DefaultUserRole)},
			mockedRes: mockedRes{isOwnerErr: someError},
			wantErr:   someError,
		},
		{
			name:    "unknown role provided",
			args:    args{userRole: "kek"},
			wantErr: NotEnoughRightsErr,
		},
		{
			name: "success",
			args: args{userRole: string(domain.AdminRole)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanSetup := setup(tt.mockedRes)
			defer cleanSetup()
			err := (&ArticleService{repoMock}).checkRights(tt.args.userOId, tt.args.userRole, tt.args.articleOId)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			assert.Nil(t, err)
		})
	}
}

func TestArticleService_UpdateReaction(t *testing.T) {
	type args struct {
		raterOId   string
		articleOId string
		reaction   string
	}
	type mockedRes struct {
		articleErr       error
		isOwnerErr       error
		oldReaction      string
		oldReactionErr   error
		reactionToUpdate string
		countToAdd       int
		updateErr        error
	}
	repoMock := repoMocks.NewArticleRepository(t)
	setup := func(res mockedRes) func() {
		articleCall := repoMock.EXPECT().GetArticle("").Return(domain.Article{}, res.articleErr).Once()
		isOwnerCall := repoMock.EXPECT().IsOwner(mock.Anything, mock.Anything).NotBefore(articleCall).
			Return(res.isOwnerErr).Maybe()
		currentReactionCall := repoMock.EXPECT().
			GetCurrentReaction("", "").NotBefore(isOwnerCall).Return(res.oldReaction, res.oldReactionErr)
		updateCall := repoMock.EXPECT().UpdateReaction(mock.Anything, mock.Anything, res.reactionToUpdate, res.countToAdd).
			NotBefore(currentReactionCall).Return(res.updateErr)
		return func() {
			articleCall.Unset()
			isOwnerCall.Unset()
			currentReactionCall.Unset()
			updateCall.Unset()
		}
	}
	someErr := errors.New("some err")
	someReaction := "🇫🇮"
	tests := []struct {
		name      string
		args      args
		mockedRes mockedRes
		wantErr   error
	}{
		{
			name:    "user is owner of article",
			wantErr: NotEnoughRightsErr,
		},
		{
			name:      "failed to execute IsOwner call",
			mockedRes: mockedRes{isOwnerErr: someErr},
			wantErr:   someErr,
		},
		{
			name:      "invalid reaction provided",
			args:      args{reaction: " "},
			mockedRes: mockedRes{isOwnerErr: NotFoundErr},
			wantErr:   InvalidDataErr,
		},
		{
			name: "failed to get current reaction",
			mockedRes: mockedRes{
				isOwnerErr:     NotFoundErr,
				oldReactionErr: someErr,
			},
			wantErr: someErr,
		},
		{
			name: "failed to remove current reaction",
			mockedRes: mockedRes{
				isOwnerErr:       NotFoundErr,
				oldReaction:      someReaction,
				reactionToUpdate: someReaction,
				countToAdd:       -1,
				updateErr:        someErr,
			},
			wantErr: someErr,
		},
		{
			name: "failed to set new reaction",
			args: args{reaction: someReaction},
			mockedRes: mockedRes{
				isOwnerErr:       NotFoundErr,
				reactionToUpdate: someReaction,
				countToAdd:       1,
				updateErr:        someErr,
			},
			wantErr: someErr,
		},
		{
			name: "ampty reaction provided",
			mockedRes: mockedRes{
				isOwnerErr: NotFoundErr,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanSetup := setup(tt.mockedRes)
			defer cleanSetup()

			err := NewArticleModel(repoMock).UpdateReaction(tt.args.raterOId, tt.args.articleOId, tt.args.reaction)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			assert.Nil(t, err)
		})
	}
}
