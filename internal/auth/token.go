package auth

import (
	"encoding/base64"
	"fmt"
	"user-management/internal/domain"
)

type Token struct {
}

func (t Token) Generate(user domain.User) string {
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", user.Username, user.Password)))
}

func (t Token) Decode(token string) (domain.User, error) {
	return domain.NewUser("", ""), nil
}
