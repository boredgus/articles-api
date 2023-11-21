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

var passwordRules = []string{"[a-z]", "[A-Z]", "[0-9]", "[./_*]"}

func (u User) Validate() error {
	logrus.Infoln(u)
	validate := validator.New()
	validate.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		for _, rule := range passwordRules {
			if !regexp.MustCompile(rule).Match([]byte(fl.Field().String())) {
				return false
			}
		}
		return true
	})
	err := validate.Struct(u)

	return parseError(err)
}

var fieldRequirements = map[string]string{
	"Password": "password should have lenth between 8 and 20, at least one lowercase letter, at least one uppercase letter, at least one number, at least one of special symbols /._*",
	"Username": "username should have length between 4 and 20",
}

func parseError(err error) (e error) {
	if err == nil {
		return nil
	}
	for _, err := range err.(validator.ValidationErrors) {
		e = fmt.Errorf("%v;", fieldRequirements[err.Field()])
	}
	return
}
