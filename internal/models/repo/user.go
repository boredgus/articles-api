package repo

import "user-management/internal/domain"

type User struct {
	OId      string          `sql:"o_id"`
	Username string          `sql:"username"`
	Password string          `sql:"pswd"`
	Role     domain.UserRole `sql:"role"`
}

type UserRepository interface {
	Create(user User) error
	Get(username string) (User, error)
	GetByOId(oid string) (User, error)
	Delete(oid string) error
	UpdateRole(oid string, role domain.UserRole) error
}
