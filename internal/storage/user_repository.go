package storage

import (
	"github.com/hop-/gotchat/internal/core"
)

type UserRepository struct {
	StorageDb
}

func newUserRepository(storage StorageDb) *UserRepository {
	return &UserRepository{storage}
}

func (r *UserRepository) GetOne(id int) (*core.User, error) {
	row := r.Db().QueryRow("SELECT id, unique_id, name, last_login FROM users WHERE id = ?", id)
	if row == nil {
		return nil, nil
	}

	var u core.User

	err := row.Scan(&u.Id, &u.UniqueId, &u.Name, &u.LastLogin)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *UserRepository) GetAll() ([]*core.User, error) {
	rows, err := r.Db().Query("SELECT id, unique_id, name, last_login FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*core.User
	for rows.Next() {
		var u core.User
		err := rows.Scan(&u.Id, &u.UniqueId, &u.Name, &u.LastLogin)
		if err != nil {
			return nil, err
		}
		users = append(users, &u)
	}

	return users, nil
}

func (r *UserRepository) GetAllBy(field string, value any) ([]*core.User, error) {
	if !isFieldExist[core.User](field) {
		return nil, ErrFieldNotExist
	}

	rows, err := r.Db().Query("SELECT id, unique_id, name, last_login FROM users where ? = ?", field, value)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*core.User
	for rows.Next() {
		var u core.User
		err := rows.Scan(&u.Id, &u.UniqueId, &u.Name, &u.LastLogin)
		if err != nil {
			return nil, err
		}
		users = append(users, &u)
	}

	return users, nil
}

func (r *UserRepository) Create(user *core.User) error {
	_, err := r.Db().Exec(
		"INSERT INTO users (unique_id, name, last_login) VALUES (?, ?, ?)",
		user.UniqueId,
		user.Name,
		user.LastLogin,
	)

	return err
}

func (r *UserRepository) Update(user *core.User) error {
	_, err := r.Db().Exec("UPDATE users SET name = ?, last_login = ? WHERE id = ?", user.Name, user.LastLogin, user.Id)

	return err
}

func (r *UserRepository) Delete(id int) error {
	_, err := r.Db().Exec("DELETE FROM users WHERE id = ?", id)

	return err
}
