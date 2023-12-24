package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArticleReaction_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		reaction ArticleReaction
		wantErr  bool
	}{
		{
			name:     "empty reaction",
			reaction: ArticleReaction(""),
			wantErr:  false,
		},
		{
			name:     "emoji with 1 unicode",
			reaction: ArticleReaction("©"),
			wantErr:  false,
		},
		{
			name:     "emoji with many unicodes",
			reaction: ArticleReaction("🏴󠁧󠁢󠁳󠁣󠁴󠁿"),
			wantErr:  false,
		},
		{
			name:     "more than 1 grapheme",
			reaction: ArticleReaction("as"),
			wantErr:  true,
		},
		{
			name:     "letter",
			reaction: ArticleReaction("w"),
			wantErr:  true,
		},
		{
			name:     "number",
			reaction: ArticleReaction("0"),
			wantErr:  true,
		},
		{
			name:     "sign",
			reaction: ArticleReaction(")"),
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.reaction.IsValid()
			assert.Equal(t, got != nil, tt.wantErr)
		})
	}
}
