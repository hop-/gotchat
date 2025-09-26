package storage

import (
	"github.com/hop-/gotchat/internal/core"
)

type MessageRepository struct {
	Repository[Message, core.Message]
}

func newMessageRepository(storage StorageDb) *MessageRepository {
	return &MessageRepository{newRepository[Message](storage)}
}

func (r *MessageRepository) Init() error {
	return r.Repository.Init()
}

func (r *MessageRepository) GetOne(id uint) (*core.Message, error) {
	m, err := r.getOne(id)
	if err != nil {
		return nil, err
	}

	return m.ToEntity(), nil
}

func (r *MessageRepository) GetOneBy(field string, value any) (*core.Message, error) {
	if !core.IsFieldExist[core.Message](field) {
		return nil, core.ErrEntityFieldNotExist
	}

	m, err := r.getOneBy(field, value)
	if err != nil {
		return nil, err
	}

	return m.ToEntity(), nil
}

func (r *MessageRepository) GetAll() ([]*core.Message, error) {
	ms, err := r.getAll()
	if err != nil {
		return nil, err
	}

	var messages []*core.Message
	for _, m := range ms {
		messages = append(messages, m.ToEntity())
	}

	return messages, nil
}

func (r *MessageRepository) GetAllBy(field string, value any) ([]*core.Message, error) {
	if !core.IsFieldExist[core.Message](field) {
		return nil, core.ErrEntityFieldNotExist
	}

	ms, err := r.getAllBy(field, value)
	if err != nil {
		return nil, err
	}

	var messages []*core.Message
	for _, m := range ms {
		messages = append(messages, m.ToEntity())
	}

	return messages, nil
}

func (r *MessageRepository) Create(entity *core.Message) (*core.Message, error) {
	m := FromEntityToMessage(entity)
	_, err := r.create(m)
	if err != nil {
		return nil, err
	}

	return m.ToEntity(), nil
}

func (r *MessageRepository) Update(entity *core.Message) (*core.Message, error) {
	m := FromEntityToMessage(entity)
	_, err := r.update(m.Id, m)
	if err != nil {
		return nil, err
	}

	return m.ToEntity(), nil
}

func (r *MessageRepository) Delete(id uint) error {
	return r.delete(id)
}
