package auth

import (
	"golang.org/x/crypto/bcrypt"
)

func NewPassword() Password {
	return Password{}
}

type Password struct{}

func (p Password) Hash(str string) (string, error) {
	res, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

func (p Password) Compare(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
