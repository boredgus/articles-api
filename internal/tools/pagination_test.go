package tools

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_pagination_Parse(t *testing.T) {
	type args struct {
		page  string
		limit string
	}
	tests := []struct {
		name    string
		args    args
		wantP   int
		wantL   int
		wantErr error
	}{
		{
			name:  "success: values are set to default",
			args:  args{},
			wantP: 0,
			wantL: 10,
		},
		{
			name:  "success: values are set to default",
			args:  args{page: "2", limit: "8"},
			wantP: 2,
			wantL: 8,
		},
		{
			name:    "page is not a number",
			args:    args{page: "q"},
			wantErr: PageNotANumberErr,
		},
		{
			name:    "page is out of range",
			args:    args{page: "-2"},
			wantErr: PageOutOfRangeErr,
		},
		{
			name:    "limit is not a number",
			args:    args{limit: "q"},
			wantErr: LimitNotANumberErr,
		},
		{
			name:    "limit is less than 0",
			args:    args{limit: "-1"},
			wantErr: LimitOutOfRangeErr,
		},
		{
			name:    "limit is greater than 50",
			args:    args{limit: "51"},
			wantErr: LimitOutOfRangeErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotP, gotL, err := NewPagination().Parse(tt.args.page, tt.args.limit)
			assert.Equal(t, gotP, tt.wantP, "check_page")
			assert.Equal(t, gotL, tt.wantL, "check_limit")
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			assert.Nil(t, err)
		})
	}
}
