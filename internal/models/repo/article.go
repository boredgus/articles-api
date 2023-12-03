package repo

import (
	"time"
	"user-management/internal/domain"
)

type ArticleData struct {
	OId   string   `sql:"o_id"`
	Theme string   `sql:"theme"`
	Text  string   `sql:"text"`
	Tags  []string `sql:"tags"`
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
	Create(userOId string, article ArticleData) error
	Get(articleOId string) (domain.Article, error)
	GetForUser(username string, page, limit int) ([]domain.Article, error)
	Update(article ArticleData) (time.Time, error)
	IsOwner(articleOId, username string) error
}
