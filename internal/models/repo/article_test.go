package repo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArticleData_CompareTags(t *testing.T) {
	type fields struct {
		Tags []string
	}
	type args struct {
		tags []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantOld []string
		wantNew []string
	}{
		{
			name:    "all tags are new",
			fields:  fields{Tags: []string{}},
			args:    args{tags: []string{"new1", "new2", "new3"}},
			wantNew: []string{"new1", "new2", "new3"},
		},
		{
			name:    "all tags are old",
			fields:  fields{Tags: []string{"old1", "old2", "old3"}},
			args:    args{tags: []string{}},
			wantOld: []string{"old1", "old2", "old3"},
		},
		{
			name:    "mixed tags",
			fields:  fields{Tags: []string{"old1", "old2", "old3"}},
			args:    args{tags: []string{"old2", "old3", "new1"}},
			wantOld: []string{"old1"},
			wantNew: []string{"new1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOld, gotNew := ArticleData{Tags: tt.fields.Tags}.CompareTags(tt.args.tags)
			assert.Equal(t, gotOld, tt.wantOld)
			assert.Equal(t, gotNew, tt.wantNew)
		})
	}
}
