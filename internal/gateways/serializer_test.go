package gateways

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJSONSerializer_Serialize(t *testing.T) {
	type args struct {
		data []any
	}
	type str struct {
		C1 string    `json:"c1"`
		C2 int       `json:"c2"`
		C3 *int      `json:"c3"`
		C4 time.Time `json:"c4"`
	}
	numb, time := 2, time.Unix(320000, 0).UTC()
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr error
	}{
		{
			name: "empty array",
			args: args{data: []any{}},
			want: "[]",
		},
		{
			name: "array of strings",
			args: args{data: []any{"1", "2", "3"}},
			want: "[\"1\",\"2\",\"3\"]",
		},
		{
			name: "array of structs",
			args: args{data: []any{
				str{C1: "1", C2: 1, C4: time},
				str{C3: &numb}}},
			want: "[{\"c1\":\"1\",\"c2\":1,\"c3\":null,\"c4\":\"1970-01-04T16:53:20Z\"},{\"c1\":\"\",\"c2\":0,\"c3\":2,\"c4\":\"0001-01-01T00:00:00Z\"}]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotE := NewJSONSerializer[any]().Serialize(tt.args.data)
			assert.Equal(t, tt.want, got)
			if tt.wantErr != nil {
				assert.ErrorIs(t, gotE, tt.wantErr)
				return
			}
			assert.Nil(t, gotE)
		})
	}
}

func TestJSONSerializer_Deserialize(t *testing.T) {
	type args struct {
		data string
	}
	type str struct {
		C1 string    `json:"c1"`
		C2 int       `json:"c2"`
		C3 *int      `json:"c3"`
		C4 time.Time `json:"c4"`
	}

	numb, tm := 2, time.Unix(320000, 0).UTC()
	tests := []struct {
		name    string
		args    args
		want    []str
		wantErr error
	}{
		{
			name: "empty array",
			args: args{data: "[]"},
			want: []str{},
		}, {
			name: "array of structs",
			args: args{data: "[{\"c1\":\"1\",\"c2\":1,\"c3\":null,\"c4\":\"1970-01-04T16:53:20Z\"},{\"c1\":\"\",\"c2\":0,\"c3\":2,\"c4\":\"0001-01-01T00:00:00Z\"}]"},
			want: []str{
				{C1: "1", C2: 1, C4: tm},
				{C3: &numb}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotE := NewJSONSerializer[str]().Deserialize(tt.args.data)
			assert.Equal(t, tt.want, got)
			if tt.wantErr != nil {
				assert.ErrorIs(t, gotE, tt.wantErr)
				return
			}
			assert.Nil(t, gotE)
		})
	}
}
