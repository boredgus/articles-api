package gateways

import (
	"strings"
	"user-management/internal/models"
)

func NewUserRepository(store Store) models.UserRepository {
	return UserRepository{store: store}
}

type UserRepository struct {
	store Store
}

func (r UserRepository) Create(user models.User) error {
	_, err := r.store.Query(`
		insert into user (o_id, username, pswd)
		values (?, ?, ?);`,
		user.OId, user.Username, user.Password)
	if err != nil && strings.Contains(err.Error(), "Error 1062") {
		return models.UsernameDuplicationErr
	}

	return err
}

func (r UserRepository) Get(username string) (models.User, error) {
	var user models.User
	rows, err := r.store.Query(`
		select o_id, username, pswd
		from user
		where user.username=?;`, username)
	if err != nil {
		return user, err
	}

	exists := rows.Next()
	if !exists {
		return user, models.InvalidAuthParameterErr
	}
	err = rows.Scan(&user.OId, &user.Username, &user.Password)
	if err != nil {
		return user, err
	}
	return user, nil
}
