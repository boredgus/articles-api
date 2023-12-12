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
	Update(userOId, userRole string, article *domain.Article) error
	Delete(userOId, userRole, articleOId string) error
}

var InvalidArticleErr = errors.New("invalid article")
var NotEnoughRightsErr = errors.New("the user does not have enough rights to perform the action")
var ArticleNotFoundErr = errors.New("article is not found")
var UnknownUserErr = errors.New("unknown user")

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
	if err == ArticleNotFoundErr {
		return UnknownUserErr
	}
	if err != nil {
		return err
	}
	article.OId = id
	article.CreatedAt = time.Now().UTC()
	return nil
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

func (a ArticleService) checkRights(userOId, userRole, articleOId string) error {
	if userRole == domain.UserRoles[domain.DefaultUserRole] {
		err := a.repo.IsOwner(articleOId, userOId)
		if err == ArticleNotFoundErr {
			return fmt.Errorf("%w: user is not an owner", NotEnoughRightsErr)
		}
		if err != nil {
			return err
		}
	} else if domain.RoleToValue[userRole] == domain.DefaultUserRole {
		return fmt.Errorf("%w: unknown user role", NotEnoughRightsErr)
	}
	return nil
}

func (a ArticleService) Update(userOId, userRole string, article *domain.Article) error {
	oldArticle, err := a.repo.Get(article.OId)
	if err != nil {
		return err
	}
	err = a.checkRights(userOId, userRole, article.OId)
	if err != nil {
		return err
	}
	if err := article.Validate(); err != nil {
		return fmt.Errorf("%w: %w", InvalidArticleErr, err)
	}
	err = a.repo.UpdateArticle(article.OId, article.Theme, article.Text)
	if err != nil {
		return err
	}
	tagsToRemove, tagsToAdd := repo.ArticleData{Tags: oldArticle.Tags}.CompareTags(article.Tags)
	if len(tagsToRemove) > 0 {
		err = a.repo.RemoveTagsFromArticle(article.OId, tagsToRemove)
		if err != nil {
			return err
		}
	}
	if len(tagsToAdd) > 0 {
		err = a.repo.AddTagsForArticle(article.OId, tagsToAdd)
		if err != nil {
			return err
		}
	}
	t := time.Now().UTC()
	article.Status = domain.UpdatedStatus
	article.CreatedAt = oldArticle.CreatedAt
	article.UpdatedAt = &t
	return nil
}

func (a ArticleService) Delete(userOId, userRole, articleOId string) error {
	article, err := a.repo.Get(articleOId)
	if err != nil {
		return err
	}
	err = a.checkRights(userOId, userRole, articleOId)
	if err != nil {
		return err
	}
	return a.repo.DeleteArticle(articleOId, article.Tags)
}
