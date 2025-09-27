package storage

import (
	"github.com/hop-/gotchat/internal/core"
)

type UserRepository struct {
	Repository[User, core.User]
}

func newUserRepository(storage StorageDb) *UserRepository {
	return &UserRepository{
		newRepository[User, core.User](storage),
	}
}

func (r *UserRepository) Init() error {
	return r.Repository.Init()
}

func (r *UserRepository) GetOne(id uint) (*core.User, error) {
	m, err := r.getOne(id)
	if err != nil {
		return nil, err
	}

	return m.ToEntity(), nil
}

func (r *UserRepository) GetOneBy(field string, value any) (*core.User, error) {
	if !core.IsFieldExist[core.User](field) {
		return nil, core.ErrEntityFieldNotExist
	}

	m, err := r.getOneBy(field, value)
	if err != nil {
		return nil, err
	}

	return m.ToEntity(), nil
}

func (r *UserRepository) GetAll() ([]*core.User, error) {
	ms, err := r.getAll()
	if err != nil {
		return nil, err
	}

	var messages []*core.User
	for _, m := range ms {
		messages = append(messages, m.ToEntity())
	}

	return messages, nil
}

func (r *UserRepository) GetAllBy(field string, value any) ([]*core.User, error) {
	if !core.IsFieldExist[core.User](field) {
		return nil, core.ErrEntityFieldNotExist
	}

	ms, err := r.getAllBy(field, value)
	if err != nil {
		return nil, err
	}

	var messages []*core.User
	for _, m := range ms {
		messages = append(messages, m.ToEntity())
	}

	return messages, nil
}

func (r *UserRepository) Create(entity *core.User) (*core.User, error) {
	m := FromEntityToUser(entity)
	_, err := r.create(m)
	if err != nil {
		return nil, err
	}

	return m.ToEntity(), nil
}

func (r *UserRepository) Update(entity *core.User) (*core.User, error) {
	m := FromEntityToUser(entity)
	_, err := r.update(m.Id, m)
	if err != nil {
		return nil, err
	}

	return m.ToEntity(), nil
}

func (r *UserRepository) Delete(id uint) error {
	return r.delete(id)
}
