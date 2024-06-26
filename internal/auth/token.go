package auth

import (
	"a-article/internal/domain"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
)

type Token[T any] interface {
	Generate(data T) (string, error)
	Decode(token string) (T, error)
}

func NewToken() Token[domain.User] {
	return BasicToken{}
}

type BasicToken struct {
}

var InvalidToken = errors.New("invalid token provided")

func (t BasicToken) Generate(user domain.User) (string, error) {
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", user.Username, user.Password))), nil
}

func (t BasicToken) Decode(token string) (u domain.User, e error) {
	bytes, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return u, fmt.Errorf("%w: %w", InvalidToken, err)
	}
	before, after, found := strings.Cut(string(bytes), ":")
	if !found {
		return u, InvalidToken
	}
	return domain.NewUser(before, after), nil
}
