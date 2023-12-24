package views

import (
	"time"
	"user-management/internal/domain"
)

// swagger:model
type Article struct {
	// identifier of article
	OId string `json:"id"`
	// theme/topic of article
	Theme string `json:"theme" form:"theme" validate:"required,min=1,max=200"`
	// content of article
	Text string `json:"text" form:"text" validate:"max=500"`
	// tags related for article
	Tags []string `json:"tags" form:"tags" validate:"tags"`
	// time of creation
	CreatedAt time.Time `json:"created_at"`
	// time when was updated
	UpdatedAt *time.Time `json:"updated_at"`
	// status of article
	// enum: deleted,created,updated
	Status string `json:"status"`
	// reactions given by users
	Reactions domain.ArticleReactions `json:"reactions"`
}

func NewArticleView(a domain.Article) Article {
	return Article{
		OId:       a.OId,
		Theme:     a.Theme,
		Text:      a.Text,
		Tags:      a.Tags,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
		Status:    a.Status.String(),
		Reactions: a.Reactions,
	}
}
