package gateways

import (
	"user-management/internal/domain"
	"user-management/internal/models"

	"github.com/google/uuid"
)

func NewUserRepository(store Store) models.UserRepository {
	return UserRepository{store: store}
}

type UserRepository struct {
	store Store
}

func (r UserRepository) Create(user domain.User) error {
	_, err := r.store.Query(`
		insert into user (o_id, username, pswd)
		values (?, ?, ?);`,
		uuid.New().String(), user.Username, user.Password)

	return err
}

func (r UserRepository) Get(username string) (domain.User, error) {
	var user domain.User
	rows, err := r.store.Query(`
		select username, pswd
		from user
		where user.username=?;`, username)
	if err != nil {
		return user, err
	}

	exists := rows.Next()
	if !exists {
		return user, models.InvalidAuthParameter
	}
	err = rows.Scan(&user.Username, &user.Password)
	if err != nil {
		return user, err
	}
	return user, nil
}
