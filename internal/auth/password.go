package auth

import (
	"golang.org/x/crypto/bcrypt"
)

func NewPassword() Pswd {
	return Pswd{}
}

type Password interface {
	Hash(str string) (string, error)
	Compare(hash, password string) bool
}

type Pswd struct{}

func (p Pswd) Hash(str string) (string, error) {
	res, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

func (p Pswd) Compare(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
