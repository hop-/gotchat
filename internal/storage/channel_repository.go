package storage

import (
	"github.com/hop-/gotchat/internal/core"
)

type ChannelRepository struct {
	Repository[Channel, core.Channel]
}

func newChannelRepository(storage StorageDb) *ChannelRepository {
	return &ChannelRepository{newRepository[Channel](storage)}
}

func (r *ChannelRepository) Init() error {
	return r.Repository.Init()
}

func (r *ChannelRepository) GetOne(id uint) (*core.Channel, error) {
	m, err := r.getOne(id)
	if err != nil {
		return nil, err
	}

	return m.ToEntity(), nil
}

func (r *ChannelRepository) GetOneBy(field string, value any) (*core.Channel, error) {
	if !core.IsFieldExist[core.Channel](field) {
		return nil, core.ErrEntityFieldNotExist
	}

	m, err := r.getOneBy(field, value)
	if err != nil {
		return nil, err
	}

	return m.ToEntity(), nil
}

func (r *ChannelRepository) GetAll() ([]*core.Channel, error) {
	ms, err := r.getAll()
	if err != nil {
		return nil, err
	}

	var messages []*core.Channel
	for _, m := range ms {
		messages = append(messages, m.ToEntity())
	}

	return messages, nil
}

func (r *ChannelRepository) GetAllBy(field string, value any) ([]*core.Channel, error) {
	if !core.IsFieldExist[core.Channel](field) {
		return nil, core.ErrEntityFieldNotExist
	}

	ms, err := r.getAllBy(field, value)
	if err != nil {
		return nil, err
	}

	var messages []*core.Channel
	for _, m := range ms {
		messages = append(messages, m.ToEntity())
	}

	return messages, nil
}

func (r *ChannelRepository) Create(entity *core.Channel) (*core.Channel, error) {
	m := FromEntityToChannel(entity)
	_, err := r.create(m)
	if err != nil {
		return nil, err
	}

	return m.ToEntity(), nil
}

func (r *ChannelRepository) Update(entity *core.Channel) (*core.Channel, error) {
	m := FromEntityToChannel(entity)
	_, err := r.update(m.Id, m)
	if err != nil {
		return nil, err
	}

	return m.ToEntity(), nil
}

func (r *ChannelRepository) Delete(id uint) error {
	return r.delete(id)
}
