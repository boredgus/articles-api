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
	UpdateReaction(raterOId, articleOId, reaction string) error
}

var NotFoundErr = errors.New("not found")
var InvalidDataErr = errors.New("invalid data")
var NotEnoughRightsErr = errors.New("the user does not have enough rights to perform the action")
var UnknownUserErr = errors.New("unknown user")

func NewArticleModel(repo repo.ArticleRepository) ArticleModel {
	return &ArticleService{repo: repo}
}

type ArticleService struct {
	repo repo.ArticleRepository
}

func (a *ArticleService) Create(userOId string, article *domain.Article) error {
	if err := article.Validate(); err != nil {
		return fmt.Errorf("%w: %w", InvalidDataErr, err)
	}
	id := uuid.New().String()
	err := a.repo.CreateArticle(userOId, repo.ArticleData{
		OId:   id,
		Theme: article.Theme,
		Text:  article.Text,
		Tags:  article.Tags,
	})
	if err != nil {
		return err
	}
	article.OId = id
	article.CreatedAt = time.Now().UTC()
	return nil
}

func (a *ArticleService) Get(articleOId string) (domain.Article, error) {
	ar, err := a.repo.GetArticle(articleOId)
	if err != nil {
		return domain.Article{}, err
	}
	reactions, err := a.repo.GetReactionsFor(articleOId)
	if err != nil {
		return domain.Article{}, err
	}
	ar.Reactions = reactions[articleOId]
	return ar, nil
}

func (a *ArticleService) GetForUser(username string, page, limit int) ([]domain.Article, PaginationData, error) {
	articles, err := a.repo.GetForUser(username, page, limit)
	if err != nil {
		return nil, PaginationData{}, err
	}
	pagination := PaginationData{Page: page, Limit: limit, Count: len(articles)}
	if len(articles) == 0 {
		return articles, pagination, nil
	}
	articledOIds := make([]string, 0, len(articles))
	for _, a := range articles {
		articledOIds = append(articledOIds, a.OId)
	}
	reactions, err := a.repo.GetReactionsFor(articledOIds...)
	if err != nil {
		return nil, PaginationData{}, err
	}
	for i, a := range articles {
		if len(reactions[a.OId]) > 0 {
			articles[i].Reactions = reactions[a.OId]
		}
	}
	return articles, pagination, nil
}

func (a *ArticleService) checkRights(userOId, userRole, articleOId string) error {
	switch domain.UserRole(userRole) {
	case domain.DefaultUserRole:
		err := a.repo.IsOwner(articleOId, userOId)
		if err == NotFoundErr {
			return fmt.Errorf("%w: user is not an owner", NotEnoughRightsErr)
		}
		return err
	case domain.ModeratorRole, domain.AdminRole:
		return nil
	default:
		return fmt.Errorf("%w: unknown user role", NotEnoughRightsErr)
	}
}

func (a *ArticleService) Update(userOId, userRole string, article *domain.Article) error {
	oldArticle, err := a.repo.GetArticle(article.OId)
	if err != nil {
		return err
	}
	err = a.checkRights(userOId, userRole, article.OId)
	if err != nil {
		return err
	}
	if err := article.Validate(); err != nil {
		return fmt.Errorf("%w: %w", InvalidDataErr, err)
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

func (a *ArticleService) Delete(userOId, userRole, articleOId string) error {
	article, err := a.repo.GetArticle(articleOId)
	if err != nil {
		return err
	}
	err = a.checkRights(userOId, userRole, articleOId)
	if err != nil {
		return err
	}
	return a.repo.DeleteArticle(articleOId, article.Tags)
}

func (a *ArticleService) UpdateReaction(raterOId, articleOId, reaction string) error {
	if _, err := a.repo.GetArticle(articleOId); err != nil {
		return err
	}
	err := a.repo.IsOwner(articleOId, raterOId)
	if err == nil {
		return fmt.Errorf("%w: it is prohibited to give reaction for own article", NotEnoughRightsErr)
	}
	if err != NotFoundErr {
		return err
	}
	if err = domain.ArticleReaction(reaction).IsValid(); err != nil {
		return fmt.Errorf("%w: %w", InvalidDataErr, err)
	}
	oldReaction, err := a.repo.GetCurrentReaction(raterOId, articleOId)
	if err != nil && err != NotFoundErr {
		return err
	}
	if len(oldReaction) > 0 {
		if err := a.repo.UpdateReaction(raterOId, articleOId, oldReaction, -1); err != nil {
			return err
		}
	}
	if reaction != string(domain.NoReaction) {
		return a.repo.UpdateReaction(raterOId, articleOId, reaction, 1)
	}
	return nil
}
