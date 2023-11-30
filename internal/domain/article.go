package domain

import (
	"fmt"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type ArticleStatus uint8

const (
	InitialStatus ArticleStatus = 0
	UpdatedStatus ArticleStatus = 1
	DeletedStatus ArticleStatus = 2
)

// swagger:model
type Article struct {
	// identifier of article
	OId string `json:"id"`
	// theme/topic of article
	Theme string `json:"theme" form:"theme" validate:"required,min=1,max=200"`
	// content of article
	Text string `json:"text" form:"text" form:"theme" validate:"max=500"`
	// tags related for article
	Tags []string `json:"tags" form:"tags" validate:"tags"`
	// time of creation
	CreatedAt time.Time `json:"created_at"`
	// time when was updated
	UpdatedAt *time.Time `json:"updated_at"`
	// status of article
	Status ArticleStatus `json:"status,omitempty"`
}

var articleRequirements = Requirements{
	"Theme": "theme should have length from 1 to 200",
	"Text":  "text should have length less than of equal 500",
	"Tags":  "tag cannot have spaces",
}

var tagsRule = regexp.MustCompile("[\t\n\f\r ]")

func (a Article) Validate() error {
	validate := validator.New()
	err := validate.RegisterValidation("tags", func(fl validator.FieldLevel) bool {
		for i := 0; i < fl.Field().Len(); i++ {
			if tagsRule.Match([]byte(fmt.Sprint(fl.Field().Index(i).Interface()))) {
				return false
			}
		}
		return true
	})
	if err != nil {
		logrus.Warnln("failed to register custom tag validation")
	}
	return parseError(validate.Struct(a), articleRequirements)
}
