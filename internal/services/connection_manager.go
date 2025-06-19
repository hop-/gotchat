package services

import (
	"context"
	"sync"

	"github.com/hop-/gotchat/internal/core"
	"github.com/hop-/gotchat/internal/network"
)

type ConnectionManager struct {
	AtomicRunningStatus
	connections   map[string]network.Conn
	mu            sync.RWMutex
	eventEmitter  core.EventEmitter
	eventListener core.EventListener
}

func NewConnectionManager(emitter core.EventEmitter, listener core.EventListener) *ConnectionManager {
	return &ConnectionManager{
		AtomicRunningStatus{},
		make(map[string]network.Conn),
		sync.RWMutex{},
		emitter,
		listener,
	}
}

// Init implements core.Service.
func (cm *ConnectionManager) Init() error {
	return nil
}

// Name implements core.Service.
func (cm *ConnectionManager) Name() string {
	return "ConnectionManager"
}

// Run implements core.Service.
func (cm *ConnectionManager) Run(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	if cm.runningStatus {
		// TODO: Handle error
		return
	}

	cm.setRunningStatus(true)

	go cm.handleEvents(ctx, wg)

	for cm.isRunning() {
		// TODO: Handle connection events
	}
}

// Close implements core.Service.
func (cm *ConnectionManager) Close() error {
	cm.setRunningStatus(false)

	return nil
}

func (cm *ConnectionManager) handleEvents(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	for cm.isRunning() {
		e, err := cm.eventListener.Next(ctx)
		if err != nil {
			// TODO: Handle error
			continue
		}

		switch e.(type) {
		// TODO: Handle specific events
		}
	}
}
