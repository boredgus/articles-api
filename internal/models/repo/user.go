package repo

import (
	"a-article/internal/domain"
	"time"
)

type User struct {
	OId      string          `sql:"o_id"`
	Username string          `sql:"username"`
	Password string          `sql:"pswd"`
	Role     domain.UserRole `sql:"role"`
}

type SignupRequest struct {
	Email       string    `sql:"email"`
	Password    string    `sql:"pswd"`
	Passcode    string    `sql:"passcode"`
	AttemptedAt time.Time `sql:"attempted_at"`
}

type UserRepository interface {
	RegisterSignupRequest(request SignupRequest) error
	GetSignupRequest(email string) (SignupRequest, error)
	Create(user User) error
	Get(username string) (User, error)
	GetByOId(oid string) (User, error)
	Delete(oid string) error
	UpdateRole(oid string, role domain.UserRole) error
}
