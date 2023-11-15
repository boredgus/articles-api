package models

import "user-management/internal/domain"

type UserModel interface {
	Create(user domain.User) error
	Authorize()
}
