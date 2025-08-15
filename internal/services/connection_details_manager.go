package services

import (
	"github.com/hop-/gotchat/internal/core"
	"github.com/hop-/gotchat/internal/storage"
)

type ConnectionDetailsManager struct {
	storage storage.Storage
}

func NewConnectionDetailsManager(storage storage.Storage) *ConnectionDetailsManager {
	return &ConnectionDetailsManager{storage: storage}
}

func (m *ConnectionDetailsManager) GetConnectionDetails(host string, client string) (*core.ConnectionDetails, error) {
	repo := m.storage.GetConnectionDetailsRepository()

	// TODO: use better approach when available (GetOneWhere, GetOneByMany, etc.)
	detailsList, err := repo.GetAllBy("host_unique_id", host)
	if err != nil {
		return nil, err
	}

	for _, details := range detailsList {
		if details.ClientUniqueId == client {
			return details, nil
		}
	}

	return nil, storage.ErrNotFound
}
