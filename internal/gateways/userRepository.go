package gateways

import (
	"strings"
	"user-management/internal/models"
	"user-management/internal/models/repo"
)

func NewUserRepository(store Store) repo.UserRepository {
	return UserRepository{store: store}
}

type UserRepository struct {
	store Store
}

func (r UserRepository) Create(user repo.User) error {
	rows, err := r.store.Query(`call CreateUser(?,?,?);`,
		user.OId, user.Username, user.Password)
	if err != nil && strings.Contains(err.Error(), "Error 1062") {
		return models.UsernameDuplicationErr
	}
	rows.Close()
	return err
}

func (r UserRepository) Get(username string) (repo.User, error) {
	var user repo.User
	rows, err := r.store.Query(`call GetUserByUsername(?);`, username)
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
	rows.Close()
	return user, nil
}

func (r UserRepository) GetByOId(oid string) (repo.User, error) {
	var user repo.User
	rows, err := r.store.Query(`call GetUserByOId(?);`, oid)
	if err != nil {
		return user, err
	}
	exists := rows.Next()
	if !exists {
		rows.Close()
		return user, models.UserNotFoundErr
	}
	err = rows.Scan(&user.OId, &user.Username, &user.Password)
	rows.Close()
	if err != nil {
		return user, err
	}
	return user, nil
}
