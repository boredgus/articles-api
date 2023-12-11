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

type UserRole int

const (
	DefaultUserRole UserRole = iota + 0
	ModeratorRole
	AdminRole
)

var userRoles = map[UserRole]string{
	DefaultUserRole: "user",
	ModeratorRole:   "moderator",
	AdminRole:       "admin",
}
var roleToValue = map[string]UserRole{
	"user":      DefaultUserRole,
	"moderator": ModeratorRole,
	"admin":     AdminRole,
}

// user credentials
// swagger:model
type User struct {
	// unique username
	// required: true
	// example: username
	Username string `json:"username" form:"username" validate:"required,min=4,max=20"`
	// secret password
	// required: true
	// example: qweQWE123.
	Password string `json:"password" form:"password" validate:"min=8,max=20,password"`
	// role of user
	// enum: user,moderator,admin
	Role string `json:"role,omitempty" form:"role" validate:"role"`
}

type Requirements map[string]string

var userRequirements = Requirements{
	"Password": "password should have lenth between 8 and 20, at least one lowercase letter, at least one uppercase letter, at least one number, at least one of special symbols .;_*/",
	"Username": "username should have length between 4 and 20",
	"Role":     "there are only user, moderator and admin roles",
}

var passwordRules = []*regexp.Regexp{
	regexp.MustCompile("[a-z]"),
	regexp.MustCompile("[A-Z]"),
	regexp.MustCompile("[0-9]"),
	regexp.MustCompile("[./_*;]")}

func (u User) GetRole() UserRole {
	return roleToValue[u.Role]
}
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
	err = validate.RegisterValidation("role", func(fl validator.FieldLevel) bool {
		value := fl.Field().String()
		if value == "" {
			return true
		}
		for _, role := range userRoles {
			if role == value {
				return true
			}
		}
		return false
	})
	if err != nil {
		logrus.Warnf("failed to register custom role validation")
	}
	return parseError(validate.Struct(u), userRequirements)
}

func parseError(err error, requirements Requirements) error {
	if err == nil {
		return nil
	}

	msg := ""
	for _, err := range err.(validator.ValidationErrors) {
		msg += requirements[err.Field()] + " ; "
	}
	return fmt.Errorf(msg)
}
