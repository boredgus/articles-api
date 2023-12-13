package gateways

import (
	"strings"
	"user-management/internal/domain"
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
	rows, err := r.store.Query(`call CreateUser(?,?,?,?);`,
		user.OId, user.Username, user.Password, user.Role)
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
	err = rows.Scan(&user.OId, &user.Username, &user.Password, &user.Role)
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
	if !rows.Next() {
		rows.Close()
		return user, models.UserNotFoundErr
	}
	err = rows.Scan(&user.OId, &user.Username, &user.Password, &user.Role)
	rows.Close()
	if err != nil {
		return user, err
	}
	return user, nil
}

func (r UserRepository) Delete(oid string) error {
	rows, err := r.store.Query("call DeleteUser(?);", oid)
	if err != nil {
		return err
	}
	rows.Close()
	return nil
}

func (r UserRepository) UpdateRole(oid string, role domain.UserRole) error {
	rows, err := r.store.Query("call UpdateUserRole(?,?);", oid, role)
	if err != nil {
		return err
	}
	rows.Close()
	return nil
}
