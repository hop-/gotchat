package core

import (
	"context"
)

type HandlerFunc func(Event)

type EventEmitter interface {
	Emit(Event)
}

type EventManager struct {
	events chan Event
}

func NewEventManager(bufferSize int) *EventManager {
	return &EventManager{
		make(chan Event, bufferSize),
	}
}

func (m *EventManager) Emit(e Event) {
	select {
	case m.events <- e:
	default:
		// Optional: log dropped events or block
		// log.Println("EventManager: dropping event", e.Type())
	}
}

func (m *EventManager) Next(ctx context.Context) (Event, error) {
	select {
	case e := <-m.events:
		return e, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
