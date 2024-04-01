package auth

import (
	"a-article/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasicToken_Generate(t *testing.T) {
	type args struct {
		user domain.User
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr error
	}{
		{
			name:    "common user",
			args:    args{user: domain.User{Username: "user", Password: "pass"}},
			want:    "dXNlcjpwYXNz",
			wantErr: nil,
		},
		{
			name:    "empty user",
			args:    args{user: domain.User{Username: "", Password: ""}},
			want:    "Og==",
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewToken().Generate(tt.args.user)
			if err != tt.wantErr {
				t.Errorf("BasicToken.Generate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("BasicToken.Generate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBasicToken_Decode(t *testing.T) {
	type args struct {
		token string
	}
	tests := []struct {
		name    string
		tr      BasicToken
		args    args
		wantU   domain.User
		wantErr error
	}{
		{
			name:    "valid token",
			args:    args{token: "dXNlcjpwYXNz"},
			wantU:   domain.NewUser("user", "pass"),
			wantErr: nil,
		},
		{
			name:    "invalid token",
			args:    args{token: "dXNlcnBhc3M="},
			wantU:   domain.User{},
			wantErr: InvalidToken,
		},
		{
			name:    "non base64 string",
			args:    args{token: "kdfjkjdflks"},
			wantU:   domain.User{},
			wantErr: InvalidToken,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotU, err := NewToken().Decode(tt.args.token)
			if err != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.Nil(t, err)
			}
			assert.Equal(t, gotU, tt.wantU)
		})
	}
}
