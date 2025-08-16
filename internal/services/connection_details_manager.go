package services

import (
	"context"
	"sync"

	"github.com/hop-/gotchat/internal/core"
)

type ConnectionDetailsManager struct {
	eventEmitter core.EventEmitter
	repo         core.Repository[core.ConnectionDetails]
}

func NewConnectionDetailsManager(eventEmitter core.EventEmitter, connectionDetailsRepo core.Repository[core.ConnectionDetails]) *ConnectionDetailsManager {
	return &ConnectionDetailsManager{eventEmitter, connectionDetailsRepo}
}

// Init implements core.Service.
func (m *ConnectionDetailsManager) Init() error {
	// Nothing to do
	return nil
}

// MapEventToCommands implements core.Service.
func (m *ConnectionDetailsManager) MapEventToCommands(event core.Event) []core.Command {
	return nil
}

// Name implements core.Service.
func (m *ConnectionDetailsManager) Name() string {
	return "ConnectionDetailsManager"
}

// Run implements core.Service.
func (m *ConnectionDetailsManager) Run(ctx context.Context, wg *sync.WaitGroup) {
	// This service does not run any background tasks.
}

// Close implements core.Service.
func (m *ConnectionDetailsManager) Close() error {
	// Nothing to do
	return nil
}

func (m *ConnectionDetailsManager) GetConnectionDetails(host string, client string) (*core.ConnectionDetails, error) {
	// TODO: use better approach when available (GetOneWhere, GetOneByMany, etc.)
	detailsList, err := m.repo.GetAllBy("host_unique_id", host)
	if err != nil {
		if err == core.ErrEntityNotFound {
			return nil, nil
		}

		return nil, err
	}

	for _, details := range detailsList {
		if details.ClientUniqueId == client {
			return details, nil
		}
	}

	return nil, nil
}
