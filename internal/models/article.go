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
	Get(articleOId string) (domain.Article, error)
	Update(userID string, article *domain.Article) error
}

var InvalidArticleErr = errors.New("invalid article")
var UserIsNotAnOwnerErr = errors.New("user does not have such article")
var ArticleNotFoundErr = errors.New("article is not found")

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
	id := uuid.New().String()
	err := a.repo.CreateArticle(userOId, repo.ArticleData{
		OId:   id,
		Theme: article.Theme,
		Text:  article.Text,
		Tags:  article.Tags,
	})
	article.OId = id
	article.CreatedAt = time.Now().UTC()

	return err
}

func (a ArticleService) Get(articleOId string) (domain.Article, error) {
	return a.repo.Get(articleOId)
}

func (a ArticleService) GetForUser(username string, page, limit int) ([]domain.Article, PaginationData, error) {
	articles, err := a.repo.GetForUser(username, page, limit)
	if err != nil {
		return nil, PaginationData{}, err
	}
	return articles, PaginationData{Page: page, Limit: limit, Count: len(articles)}, nil
}

func (a ArticleService) Update(username string, article *domain.Article) error {
	oldArticle, err := a.repo.IsOwner(article.OId, username)
	fmt.Println("> is owner res", oldArticle, err)
	if err != nil {
		return err
	}
	if err := article.Validate(); err != nil {
		return fmt.Errorf("%w: %w", InvalidArticleErr, err)
	}
	err = a.repo.UpdateArticle(article.OId, article.Theme, article.Text)
	fmt.Println("> update err", err)
	if err != nil {
		return err
	}
	tagsToRemove, tagsToAdd := repo.ArticleData{Tags: oldArticle.Tags}.CompareTags(article.Tags)
	fmt.Println("> compare res", tagsToRemove, tagsToAdd)
	if len(tagsToRemove) > 0 {
		err = a.repo.RemoveTagsFromArticle(article.OId, tagsToRemove)
		fmt.Println("> remove tags err", err)
		if err != nil {
			return err
		}
	}
	if len(tagsToAdd) > 0 {
		err = a.repo.AddTagsForArticle(article.OId, tagsToAdd)
		fmt.Println("> add tags err", err)
		if err != nil {
			return err
		}
	}
	t := time.Now().UTC()
	article.Status = domain.UpdatedStatus
	article.CreatedAt = oldArticle.CreatedAt
	article.UpdatedAt = &t
	fmt.Println("> success")
	return nil
}
