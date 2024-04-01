package gateways

import (
	"a-article/internal/domain"
	"a-article/internal/models"
	"a-article/internal/models/repo"
	e "a-article/pkg/db/errors"
)

func NewUserRepository(store Store) repo.UserRepository {
	return &UserRepository{store: store}
}

type UserRepository struct {
	store Store
}

func (r *UserRepository) Create(user repo.User) error {
	rows, err := r.store.Query(`call articlesdb.CreateUser($1,$2,$3,$4);`,
		user.OId, user.Username, user.Password, user.Role)
	if e.IsPqErrorCode(err, e.UniqueViolationError) {
		return models.UsernameDuplicationErr
	}
	if err != nil {
		return err
	}
	return rows.Close()
}

func (r *UserRepository) Get(username string) (repo.User, error) {
	var user repo.User
	rows, err := r.store.Query(`select * from articlesdb.GetUserByUsername($1);`, username)
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

func (r *UserRepository) GetByOId(oid string) (repo.User, error) {
	var user repo.User
	rows, err := r.store.Query(`select * from articlesdb.GetUserByOId($1);`, oid)
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

func (r *UserRepository) Delete(oid string) error {
	rows, err := r.store.Query("call DeleteUser($1);", oid)
	if err != nil {
		return err
	}
	rows.Close()
	return nil
}

func (r *UserRepository) UpdateRole(oid string, role domain.UserRole) error {
	rows, err := r.store.Query("call UpdateUserRole($1,$2);", oid, role)
	if err != nil {
		return err
	}
	rows.Close()
	return nil
}
