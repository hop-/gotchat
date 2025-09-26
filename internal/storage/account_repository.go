package storage

import (
	"github.com/hop-/gotchat/internal/core"
)

type AccountRepository struct {
	Repository[Account, core.Account]
}

func newAccountRepository(storage StorageDb) *AccountRepository {
	return &AccountRepository{newRepository[Account](storage)}
}

func (r *AccountRepository) GetOne(id uint) (*core.Account, error) {
	m, err := r.getOne(id)
	if err != nil {
		return nil, err
	}

	return m.ToEntity(), nil
}

func (r *AccountRepository) GetOneBy(field string, value any) (*core.Account, error) {
	if !core.IsFieldExist[core.Account](field) {
		return nil, core.ErrEntityFieldNotExist
	}

	m, err := r.getOneBy(field, value)
	if err != nil {
		return nil, err
	}

	return m.ToEntity(), nil
}

func (r *AccountRepository) GetAll() ([]*core.Account, error) {
	ms, err := r.getAll()
	if err != nil {
		return nil, err
	}

	var messages []*core.Account
	for _, m := range ms {
		messages = append(messages, m.ToEntity())
	}

	return messages, nil
}

func (r *AccountRepository) GetAllBy(field string, value any) ([]*core.Account, error) {
	if !core.IsFieldExist[core.Account](field) {
		return nil, core.ErrEntityFieldNotExist
	}

	ms, err := r.getAllBy(field, value)
	if err != nil {
		return nil, err
	}

	var messages []*core.Account
	for _, m := range ms {
		messages = append(messages, m.ToEntity())
	}

	return messages, nil
}

func (r *AccountRepository) Create(entity *core.Account) (*core.Account, error) {
	m := FromEntityToAccount(entity)
	_, err := r.create(m)
	if err != nil {
		return nil, err
	}

	return m.ToEntity(), nil
}

func (r *AccountRepository) Update(entity *core.Account) (*core.Account, error) {
	m := FromEntityToAccount(entity)
	_, err := r.update(m.Id, m)
	if err != nil {
		return nil, err
	}

	return m.ToEntity(), nil
}

func (r *AccountRepository) Delete(id uint) error {
	return r.delete(id)
}
