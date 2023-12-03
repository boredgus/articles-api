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
	rows, err := r.store.Query(`
		insert into user (o_id, username, pswd)
		values (?, ?, ?);`,
		user.OId, user.Username, user.Password)
	rows.Close()
	if err != nil && strings.Contains(err.Error(), "Error 1062") {
		return models.UsernameDuplicationErr
	}

	return err
}

func (r UserRepository) Get(username string) (repo.User, error) {
	var user repo.User
	rows, err := r.store.Query(`
		select o_id, username, pswd
		from user
		where user.username=?;`, username)
	if err != nil {
		rows.Close()
		return user, err
	}

	exists := rows.Next()
	if !exists {
		rows.Close()
		return user, models.InvalidAuthParameterErr
	}
	err = rows.Scan(&user.OId, &user.Username, &user.Password)
	rows.Close()
	if err != nil {
		return user, err
	}
	return user, nil
}

func (r UserRepository) GetByOId(oid string) (repo.User, error) {
	var user repo.User
	rows, err := r.store.Query(`
		select o_id, username, pswd
		from user
		where user.o_id=?;`, oid)
	if err != nil {
		rows.Close()
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
