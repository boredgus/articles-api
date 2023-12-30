package server

import (
	"context"
	"errors"
	"fmt"
	pb "user-management/grpc"
	"user-management/internal/domain"
	"user-management/internal/models"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewArticleService(model models.ArticleModel) pb.ArticleServiceServer {
	return &ArticleServer{model: model}
}

type ArticleServer struct {
	pb.UnimplementedArticleServiceServer
	model models.ArticleModel
}

func (s *ArticleServer) GetArticle(ctx context.Context, req *pb.GetArticleRequest) (*pb.GetArticleReply, error) {
	a, err := s.model.Get(req.ArticleId)
	if errors.Is(err, models.NotFoundErr) {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.GetArticleReply{Article: a.ToProto()}, nil
}

func (s *ArticleServer) GetArticles(ctx context.Context, req *pb.GetArticlesRequest) (*pb.GetArticlesReply, error) {
	a, p, err := s.model.GetForUser(req.Username, int(req.Page), int(req.Limit))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	articles := make([]*pb.Article, len(a))
	for i := range a {
		articles[i] = a[i].ToProto()
	}
	return &pb.GetArticlesReply{
		Articles: articles,
		Pagination: &pb.PaginationData{
			Page:  int32(p.Page),
			Limit: int32(p.Limit),
			Count: int32(p.Count),
		},
	}, nil
}

func (s *ArticleServer) CreateArticle(ctx context.Context, req *pb.CreateArticleRequest) (*pb.CreateArticleReply, error) {
	a := domain.FromProtoData(req.Article)
	err := s.model.Create(req.UserId, a)
	if errors.Is(err, models.InvalidDataErr) {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.CreateArticleReply{Article: a.ToProto()}, nil
}

func (s *ArticleServer) UpdateArticle(ctx context.Context, req *pb.UpdateArticleRequest) (*pb.UpdateArticleReply, error) {
	a := domain.FromProtoData(req.Article)
	a.OId = req.ArticleId
	err := s.model.Update(req.UserId, req.UserRole, a)
	if errors.Is(err, models.NotFoundErr) {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	if errors.Is(err, models.NotEnoughRightsErr) {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}
	if errors.Is(err, models.InvalidDataErr) {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.UpdateArticleReply{Article: a.ToProto()}, nil
}

func (s *ArticleServer) DeleteArticle(ctx context.Context, req *pb.DeleteArticleRequest) (*pb.MessageReply, error) {
	err := s.model.Delete(req.UserId, req.UserRole, req.ArticleId)
	if errors.Is(err, models.NotFoundErr) {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	if errors.Is(err, models.NotEnoughRightsErr) {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.MessageReply{Message: "successfuly deleted"}, nil
}

func (s *ArticleServer) UpdateReaction(ctx context.Context, req *pb.UpdateReactionRequest) (*pb.MessageReply, error) {
	err := s.model.UpdateReaction(req.RaterId, req.ArticleId, req.Reaction)
	if errors.Is(err, models.InvalidDataErr) {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if errors.Is(err, models.NotFoundErr) {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	if errors.Is(err, models.NotEnoughRightsErr) {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.MessageReply{Message: fmt.Sprintf("reaction set to '%v'", req.Reaction)}, nil
}
