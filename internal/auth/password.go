package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func NewCryptor() Crptr {
	return Crptr{}
}

type Cryptor interface {
	Encrypt(str string) (string, error)
	Compare(hash, comparedStr string) bool
}

type Crptr struct{}

func (p Crptr) Encrypt(str string) (string, error) {
	res, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

func (p Crptr) Compare(hash, comparedStr string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(comparedStr))
	fmt.Printf("> comparing '%s' with '%s': %v\n\n\n", comparedStr, hash, err)
	return err == nil
}
