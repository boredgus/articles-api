package models

import (
	"errors"
	"fmt"
	"time"
	"user-management/internal/domain"
	"user-management/internal/models/repo"

	"github.com/google/uuid"
)

type PaginationData struct {
	Page  int `json:"page" sql:"page"`
	Limit int `json:"limit" sql:"limit"`
	Count int `json:"count" sql:"count"`
}

type ArticleModel interface {
	Create(userOId string, article *domain.Article) error
	GetForUser(username string, page, limit int) ([]domain.Article, PaginationData, error)
	Update(userID string, article *domain.Article) error
}

var InvalidArticleErr = errors.New("invalid article")
var UserIsNotAnOwnerErr = errors.New("user does not have such article")

func NewArticleModel(repo repo.ArticleRepository) ArticleModel {
	return ArticleService{repo}
}

type ArticleService struct {
	repo repo.ArticleRepository
}

func (a ArticleService) Create(userOId string, article *domain.Article) error {
	if err := article.Validate(); err != nil {
		return fmt.Errorf("%w: %w", InvalidArticleErr, err)
	}
	fmt.Print("> before repo call")
	id := uuid.New().String()
	err := a.repo.Create(userOId, repo.ArticleData{
		OId:   id,
		Theme: article.Theme,
		Text:  article.Text,
		Tags:  article.Tags,
	})
	article.OId = id
	article.CreatedAt = time.Now()

	return err
}

func (a ArticleService) GetForUser(username string, page, limit int) ([]domain.Article, PaginationData, error) {
	articles, err := a.repo.GetForUser(username, page, limit)
	if err != nil {
		return nil, PaginationData{}, err
	}
	return articles, PaginationData{Page: page, Limit: limit, Count: len(articles)}, nil
}

func (a ArticleService) Update(username string, article *domain.Article) error {
	err := a.repo.IsOwner(article.OId, username)
	if errors.Is(err, UserIsNotAnOwnerErr) {
		return err
	}
	if err := article.Validate(); err != nil {
		return fmt.Errorf("%w: %w", InvalidArticleErr, err)
	}
	err = a.repo.Update(repo.ArticleData{
		OId:   article.OId,
		Theme: article.Theme,
		Text:  article.Text,
		Tags:  article.Tags,
	})
	t := time.Now()
	article.UpdatedAt = &t
	return err
}
