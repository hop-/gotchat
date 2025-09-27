package storage

import "github.com/hop-/gotchat/internal/core"

type ConnectionDetailsRepository struct {
	Repository[ConnectionDetails, core.ConnectionDetails]
}

func newConnectionDetailsRepository(storage StorageDb) *ConnectionDetailsRepository {
	return &ConnectionDetailsRepository{
		newRepository[ConnectionDetails](storage),
	}
}

func (r *ConnectionDetailsRepository) Init() error {
	return r.Repository.Init()
}

func (r *ConnectionDetailsRepository) GetOne(id uint) (*core.ConnectionDetails, error) {
	m, err := r.getOne(id)
	if err != nil {
		return nil, err
	}

	return m.ToEntity(), nil
}

func (r *ConnectionDetailsRepository) GetOneBy(field string, value any) (*core.ConnectionDetails, error) {
	if !core.IsFieldExist[core.ConnectionDetails](field) {
		return nil, core.ErrEntityFieldNotExist
	}

	m, err := r.getOneBy(field, value)
	if err != nil {
		return nil, err
	}

	return m.ToEntity(), nil
}

func (r *ConnectionDetailsRepository) GetAll() ([]*core.ConnectionDetails, error) {
	ms, err := r.getAll()
	if err != nil {
		return nil, err
	}

	var messages []*core.ConnectionDetails
	for _, m := range ms {
		messages = append(messages, m.ToEntity())
	}

	return messages, nil
}

func (r *ConnectionDetailsRepository) GetAllBy(field string, value any) ([]*core.ConnectionDetails, error) {
	if !core.IsFieldExist[core.ConnectionDetails](field) {
		return nil, core.ErrEntityFieldNotExist
	}

	ms, err := r.getAllBy(field, value)
	if err != nil {
		return nil, err
	}

	var messages []*core.ConnectionDetails
	for _, m := range ms {
		messages = append(messages, m.ToEntity())
	}

	return messages, nil
}

func (r *ConnectionDetailsRepository) Create(entity *core.ConnectionDetails) (*core.ConnectionDetails, error) {
	m := FromEntityToConnectionDetails(entity)
	_, err := r.create(m)
	if err != nil {
		return nil, err
	}

	return m.ToEntity(), nil
}

func (r *ConnectionDetailsRepository) Update(entity *core.ConnectionDetails) (*core.ConnectionDetails, error) {
	m := FromEntityToConnectionDetails(entity)
	_, err := r.update(m.Id, m)
	if err != nil {
		return nil, err
	}

	return m.ToEntity(), nil
}

func (r *ConnectionDetailsRepository) Delete(id uint) error {
	return r.delete(id)
}
