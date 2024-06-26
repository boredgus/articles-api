package domain

import (
	"fmt"
	"strings"
	"time"

	grpc "a-article/grpc"

	"github.com/go-playground/validator/v10"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/rivo/uniseg"
	"github.com/sirupsen/logrus"
)

type ArticleStatus int

const (
	DeletedStatus ArticleStatus = iota - 1
	InitialStatus
	UpdatedStatus
)

var statuses = map[ArticleStatus]string{
	DeletedStatus: "deleted",
	InitialStatus: "created",
	UpdatedStatus: "updated",
}

func (s ArticleStatus) String() string {
	return statuses[s]
}

type ArticleReactions map[string]int32

type ArticleReaction string

const NoReaction ArticleReaction = ""

func (reaction ArticleReaction) IsValid() error {
	if reaction == NoReaction {
		return nil
	}
	if uniseg.GraphemeClusterCount(string(reaction)) > 1 {
		return fmt.Errorf("reaction should have only one grapheme")
	}
	if len(reaction) == 1 {
		return fmt.Errorf("reaction should be an emoji")
	}
	return nil
}

type Article struct {
	OId       string           `json:"id"`
	Theme     string           `json:"theme" form:"theme" validate:"required,min=1,max=200"`
	Text      string           `json:"text" form:"text" validate:"max=500"`
	Tags      []string         `json:"tags" form:"tags" validate:"tags"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt *time.Time       `json:"updated_at"`
	Status    ArticleStatus    `json:"status,omitempty"`
	Reactions ArticleReactions `json:"reactions"`
}

var articleRequirements = Requirements{
	"Theme": "theme should have length from 1 to 200",
	"Text":  "text should have length less than of equal 500",
	"Tags":  "tag cannot have spaces",
}

func (a *Article) Validate() error {
	validate := validator.New()
	err := validate.RegisterValidation("tags", func(fl validator.FieldLevel) bool {
		for i := 0; i < fl.Field().Len(); i++ {
			if strings.Contains(fmt.Sprint(fl.Field().Index(i).Interface()), " ") {
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

func (a *Article) ToProto() *grpc.Article {
	var updatedAt *timestamp.Timestamp
	if a.UpdatedAt != nil {
		updatedAt = &timestamp.Timestamp{Seconds: a.UpdatedAt.Unix()}
	}
	return &grpc.Article{
		Id:        a.OId,
		Theme:     a.Theme,
		Text:      a.Text,
		Tags:      a.Tags,
		Reactions: a.Reactions,
		Status:    grpc.ArticleStatus(a.Status),
		CreatedAt: &timestamp.Timestamp{Seconds: a.CreatedAt.Unix()},
		UpdatedAt: updatedAt,
	}
}

func FromProtoData(protoDTO *grpc.ArticleData) *Article {
	return &Article{
		Theme: protoDTO.Theme,
		Text:  protoDTO.Text,
		Tags:  protoDTO.Tags,
	}
}
