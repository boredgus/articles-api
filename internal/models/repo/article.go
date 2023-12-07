package repo

import (
	"slices"
	"time"
	"user-management/internal/domain"
)

type ArticleData struct {
	OId   string   `sql:"o_id"`
	Theme string   `sql:"theme"`
	Text  string   `sql:"text"`
	Tags  []string `sql:"tags"`
}

func (a ArticleData) CompareTags(tags []string) (old []string, new []string) {
	for _, curr := range a.Tags {
		if !slices.Contains[[]string, string](tags, curr) {
			old = append(old, curr)
		}
	}
	for _, comp := range tags {
		if !slices.Contains[[]string, string](a.Tags, comp) {
			new = append(new, comp)
		}
	}
	return
}

type Article struct {
	OId       string     `sql:"o_id"`
	Theme     string     `sql:"theme"`
	Text      string     `sql:"text"`
	Tags      []string   `sql:"tags"`
	CreatedAt time.Time  `sql:"created_at"`
	UpdatedAt *time.Time `sql:"updated_at"`
}

type ArticleRepository interface {
	CreateArticle(userOId string, article ArticleData) error
	Get(articleOId string) (domain.Article, error)
	GetForUser(username string, page, limit int) ([]domain.Article, error)
	IsOwner(articleOId, username string) (domain.Article, error)
	UpdateArticle(oid, theme, text string) error
	AddTagsForArticle(articleOId string, tags []string) error
	RemoveTagsFromArticle(articleOId string, tags []string) error
}
