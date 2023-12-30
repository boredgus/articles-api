package server

import (
	"context"
	"fmt"
	"testing"
	pb "user-management/grpc"
	"user-management/internal/domain"
	mdlMocks "user-management/internal/mocks/models"
	"user-management/internal/models"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestArticleServer_GetArticle(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.GetArticleRequest
	}
	type mockedRes struct {
		article    domain.Article
		articleErr error
	}
	modelMock := mdlMocks.NewArticleModel(t)
	setup := func(res *mockedRes, args *pb.GetArticleRequest) func() {
		modelCall := modelMock.EXPECT().Get(args.ArticleId).Return(res.article, res.articleErr).Once()
		return func() {
			modelCall.Unset()
		}
	}
	someErr := fmt.Errorf("some err")
	article := domain.Article{}
	tests := []struct {
		name      string
		mockedRes mockedRes
		args      args
		want      *pb.GetArticleReply
		wantErr   error
	}{
		{
			name:      "article not found",
			args:      args{req: &pb.GetArticleRequest{}},
			mockedRes: mockedRes{articleErr: models.NotFoundErr},
			wantErr:   status.Error(codes.NotFound, models.NotFoundErr.Error()),
		},
		{
			name:      "internal error",
			args:      args{req: &pb.GetArticleRequest{}},
			mockedRes: mockedRes{articleErr: someErr},
			wantErr:   status.Error(codes.Internal, someErr.Error()),
		},
		{
			name:      "success",
			args:      args{req: &pb.GetArticleRequest{}},
			mockedRes: mockedRes{article: article},
			want:      &pb.GetArticleReply{Article: article.ToProto()},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanSetup := setup(&tt.mockedRes, tt.args.req)
			defer cleanSetup()

			got, err := (&ArticleServer{
				UnimplementedArticleServiceServer: pb.UnimplementedArticleServiceServer{},
				model:                             modelMock,
			}).GetArticle(tt.args.ctx, tt.args.req)
			assert.Equal(t, got, tt.want)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			assert.Nil(t, err)
		})
	}
}

func TestArticleServer_GetArticles(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.GetArticlesRequest
	}
	type mockedRes struct {
		articles    []domain.Article
		pagination  models.PaginationData
		articlesErr error
	}
	modelMock := mdlMocks.NewArticleModel(t)
	setup := func(res *mockedRes, args *pb.GetArticlesRequest) func() {
		modelCall := modelMock.EXPECT().
			GetForUser(args.Username, int(args.Page), int(args.Limit)).Return(res.articles, res.pagination, res.articlesErr).Once()
		return func() {
			modelCall.Unset()
		}
	}
	someErr := fmt.Errorf("some err")
	articles := []domain.Article{{}}
	tests := []struct {
		name      string
		mockedRes mockedRes
		args      args
		want      *pb.GetArticlesReply
		wantErr   error
	}{
		{
			name:      "internal error",
			args:      args{req: &pb.GetArticlesRequest{}},
			mockedRes: mockedRes{articlesErr: someErr},
			wantErr:   status.Error(codes.Internal, someErr.Error()),
		},
		{
			name: "success",
			args: args{req: &pb.GetArticlesRequest{
				Username: "user",
				Page:     1,
				Limit:    2,
			}},
			mockedRes: mockedRes{
				articles:   articles,
				pagination: models.PaginationData{Page: 1, Limit: 2, Count: 1}},
			want: &pb.GetArticlesReply{
				Articles:   []*pb.Article{articles[0].ToProto()},
				Pagination: &pb.PaginationData{Page: 1, Limit: 2, Count: 1}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanSetup := setup(&tt.mockedRes, tt.args.req)
			defer cleanSetup()

			got, err := (&ArticleServer{
				UnimplementedArticleServiceServer: pb.UnimplementedArticleServiceServer{},
				model:                             modelMock,
			}).GetArticles(tt.args.ctx, tt.args.req)
			assert.Equal(t, got, tt.want)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			assert.Nil(t, err)
		})
	}
}

func TestArticleServer_CreateArticle(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.CreateArticleRequest
	}
	type mockedRes struct {
		createErr error
	}
	modelMock := mdlMocks.NewArticleModel(t)
	setup := func(res *mockedRes, args *pb.CreateArticleRequest) func() {
		modelCall := modelMock.EXPECT().
			Create(args.UserId,
				&domain.Article{Theme: args.Article.Theme, Text: args.Article.Text, Tags: args.Article.Tags},
			).Return(res.createErr).Once()
		return func() {
			modelCall.Unset()
		}
	}
	someErr := fmt.Errorf("some err")
	tests := []struct {
		name      string
		mockedRes mockedRes
		args      args
		want      *pb.CreateArticleReply
		wantErr   error
	}{
		{
			name:      "invalid article provided",
			args:      args{req: &pb.CreateArticleRequest{Article: &pb.ArticleData{}}},
			mockedRes: mockedRes{createErr: models.InvalidDataErr},
			wantErr:   status.Error(codes.InvalidArgument, models.InvalidDataErr.Error()),
		},
		{
			name:      "internal error",
			args:      args{req: &pb.CreateArticleRequest{Article: &pb.ArticleData{}}},
			mockedRes: mockedRes{createErr: someErr},
			wantErr:   status.Error(codes.Internal, someErr.Error()),
		},
		{
			name: "success",
			args: args{req: &pb.CreateArticleRequest{Article: &pb.ArticleData{}}},
			want: &pb.CreateArticleReply{Article: (&domain.Article{}).ToProto()},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanSetup := setup(&tt.mockedRes, tt.args.req)
			defer cleanSetup()

			got, err := (&ArticleServer{
				UnimplementedArticleServiceServer: pb.UnimplementedArticleServiceServer{},
				model:                             modelMock,
			}).CreateArticle(tt.args.ctx, tt.args.req)
			assert.Equal(t, got, tt.want)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			assert.Nil(t, err)
		})
	}
}

func TestArticleServer_UpdateArticle(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.UpdateArticleRequest
	}
	type mockedRes struct {
		updateErr error
	}
	modelMock := mdlMocks.NewArticleModel(t)
	setup := func(res *mockedRes, args *pb.UpdateArticleRequest) func() {
		modelCall := modelMock.EXPECT().
			Update(args.UserId, args.UserRole,
				&domain.Article{OId: args.ArticleId, Theme: args.Article.Theme, Text: args.Article.Text, Tags: args.Article.Tags},
			).Return(res.updateErr).Once()
		return func() {
			modelCall.Unset()
		}
	}
	someErr := fmt.Errorf("some err")
	tests := []struct {
		name      string
		mockedRes mockedRes
		args      args
		want      *pb.UpdateArticleReply
		wantErr   error
	}{
		{
			name:      "article not found",
			args:      args{req: &pb.UpdateArticleRequest{Article: &pb.ArticleData{}}},
			mockedRes: mockedRes{updateErr: models.NotFoundErr},
			wantErr:   status.Error(codes.NotFound, models.NotFoundErr.Error()),
		},
		{
			name:      "permission denied",
			args:      args{req: &pb.UpdateArticleRequest{Article: &pb.ArticleData{}}},
			mockedRes: mockedRes{updateErr: models.NotEnoughRightsErr},
			wantErr:   status.Error(codes.PermissionDenied, models.NotEnoughRightsErr.Error()),
		},
		{
			name:      "invalid article provided",
			args:      args{req: &pb.UpdateArticleRequest{Article: &pb.ArticleData{}}},
			mockedRes: mockedRes{updateErr: models.InvalidDataErr},
			wantErr:   status.Error(codes.InvalidArgument, models.InvalidDataErr.Error()),
		},
		{
			name:      "internal error",
			args:      args{req: &pb.UpdateArticleRequest{Article: &pb.ArticleData{}}},
			mockedRes: mockedRes{updateErr: someErr},
			wantErr:   status.Error(codes.Internal, someErr.Error()),
		},
		{
			name: "success",
			args: args{req: &pb.UpdateArticleRequest{Article: &pb.ArticleData{}}},
			want: &pb.UpdateArticleReply{Article: (&domain.Article{}).ToProto()},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanSetup := setup(&tt.mockedRes, tt.args.req)
			defer cleanSetup()

			got, err := (&ArticleServer{
				UnimplementedArticleServiceServer: pb.UnimplementedArticleServiceServer{},
				model:                             modelMock,
			}).UpdateArticle(tt.args.ctx, tt.args.req)
			assert.Equal(t, got, tt.want)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			assert.Nil(t, err)
		})
	}
}

func TestArticleServer_DeleteArticle(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.DeleteArticleRequest
	}
	type mockedRes struct {
		deleteErr error
	}
	modelMock := mdlMocks.NewArticleModel(t)
	setup := func(res *mockedRes, args *pb.DeleteArticleRequest) func() {
		modelCall := modelMock.EXPECT().
			Delete(args.UserId, args.UserRole, args.ArticleId).Return(res.deleteErr).Once()
		return func() {
			modelCall.Unset()
		}
	}
	someErr := fmt.Errorf("some err")
	tests := []struct {
		name      string
		mockedRes mockedRes
		args      args
		want      *pb.MessageReply
		wantErr   error
	}{
		{
			name:      "article not found",
			args:      args{req: &pb.DeleteArticleRequest{}},
			mockedRes: mockedRes{deleteErr: models.NotFoundErr},
			wantErr:   status.Error(codes.NotFound, models.NotFoundErr.Error()),
		},
		{
			name:      "permission denied",
			args:      args{req: &pb.DeleteArticleRequest{}},
			mockedRes: mockedRes{deleteErr: models.NotEnoughRightsErr},
			wantErr:   status.Error(codes.PermissionDenied, models.NotEnoughRightsErr.Error()),
		},
		{
			name:      "internal error",
			args:      args{req: &pb.DeleteArticleRequest{}},
			mockedRes: mockedRes{deleteErr: someErr},
			wantErr:   status.Error(codes.Internal, someErr.Error()),
		},
		{
			name: "success",
			args: args{req: &pb.DeleteArticleRequest{}},
			want: &pb.MessageReply{Message: "successfuly deleted"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanSetup := setup(&tt.mockedRes, tt.args.req)
			defer cleanSetup()

			got, err := (&ArticleServer{
				UnimplementedArticleServiceServer: pb.UnimplementedArticleServiceServer{},
				model:                             modelMock,
			}).DeleteArticle(tt.args.ctx, tt.args.req)
			assert.Equal(t, got, tt.want)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			assert.Nil(t, err)
		})
	}
}

func TestArticleServer_UpdateReaction(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.UpdateReactionRequest
	}
	type mockedRes struct {
		updateErr error
	}
	modelMock := mdlMocks.NewArticleModel(t)
	setup := func(res *mockedRes, args *pb.UpdateReactionRequest) func() {
		modelCall := modelMock.EXPECT().
			UpdateReaction(args.RaterId, args.ArticleId, args.Reaction).Return(res.updateErr).Once()
		return func() {
			modelCall.Unset()
		}
	}
	someErr := fmt.Errorf("some err")
	tests := []struct {
		name      string
		mockedRes mockedRes
		args      args
		want      *pb.MessageReply
		wantErr   error
	}{
		{
			name:      "article not found",
			args:      args{req: &pb.UpdateReactionRequest{}},
			mockedRes: mockedRes{updateErr: models.NotFoundErr},
			wantErr:   status.Error(codes.NotFound, models.NotFoundErr.Error()),
		},
		{
			name:      "permission denied",
			args:      args{req: &pb.UpdateReactionRequest{}},
			mockedRes: mockedRes{updateErr: models.NotEnoughRightsErr},
			wantErr:   status.Error(codes.PermissionDenied, models.NotEnoughRightsErr.Error()),
		},
		{
			name:      "internal error",
			args:      args{req: &pb.UpdateReactionRequest{}},
			mockedRes: mockedRes{updateErr: someErr},
			wantErr:   status.Error(codes.Internal, someErr.Error()),
		},
		{
			name: "success",
			args: args{req: &pb.UpdateReactionRequest{Reaction: "reaction"}},
			want: &pb.MessageReply{Message: "reaction set to 'reaction'"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanSetup := setup(&tt.mockedRes, tt.args.req)
			defer cleanSetup()

			got, err := (&ArticleServer{
				UnimplementedArticleServiceServer: pb.UnimplementedArticleServiceServer{},
				model:                             modelMock,
			}).UpdateReaction(tt.args.ctx, tt.args.req)
			assert.Equal(t, got, tt.want)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			assert.Nil(t, err)
		})
	}
}
