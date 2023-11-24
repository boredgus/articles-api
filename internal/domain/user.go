package domain

import (
	"fmt"
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

func NewUser(name, pswd string) User {
	return User{Username: name, Password: pswd}
}

type User struct {
	Username string `json:"username" sql:"username" form:"username" validate:"required,min=4,max=20"`
	Password string `json:"password" sql:"pswd" form:"password" validate:"min=8,max=20,password"`
}

var passwordRules = []*regexp.Regexp{
	regexp.MustCompile("[a-z]"),
	regexp.MustCompile("[A-Z]"),
	regexp.MustCompile("[0-9]"),
	regexp.MustCompile("[./_*;]")}

func (u User) Validate() error {
	validate := validator.New()
	err := validate.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		for _, rule := range passwordRules {
			if !rule.Match([]byte(fl.Field().String())) {
				return false
			}
		}
		return true
	})
	if err != nil {
		logrus.Warnf("failed to register custom password validation")
	}
	err = validate.Struct(u)

	return parseError(err)
}

var fieldRequirements = map[string]string{
	"Password": "password should have lenth between 8 and 20, at least one lowercase letter, at least one uppercase letter, at least one number, at least one of special symbols .;_*/",
	"Username": "username should have length between 4 and 20",
}

func parseError(err error) error {
	if err == nil {
		return nil
	}

	msg := ""
	for _, err := range err.(validator.ValidationErrors) {
		msg += fieldRequirements[err.Field()] + " ; "
	}
	return fmt.Errorf(msg)
}
