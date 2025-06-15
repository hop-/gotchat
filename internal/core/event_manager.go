package core

import (
	"context"
	"sync"
)

type HandlerFunc func(Event)

type EventEmitter interface {
	Emit(Event)
}

type EventListener chan Event

type EventManager struct {
	bufferSize  int
	listenersMu sync.RWMutex
	listeners   []EventListener
}

func NewEventManager(bufferSize int) *EventManager {
	return &EventManager{
		bufferSize,
		sync.RWMutex{},
		make([]EventListener, 0),
		//make(chan Event, bufferSize),
	}
}

func (m *EventManager) Register(ctx context.Context) EventListener {
	ch := make(EventListener, m.bufferSize)
	m.listenersMu.Lock()
	m.listeners = append(m.listeners, ch)
	m.listenersMu.Unlock()

	go func() {
		<-ctx.Done()
		m.Unregister(ch)
	}()

	return ch
}

func (m *EventManager) Unregister(ch EventListener) {
	m.listenersMu.Lock()
	defer m.listenersMu.Unlock()

	for i, listener := range m.listeners {
		if listener == ch {
			m.listeners = append(m.listeners[:i], m.listeners[i+1:]...)
			close(listener) // Close the channel to signal that it's no longer in use

			return
		}
	}
}

func (m *EventManager) Emit(e Event) {
	m.listenersMu.RLock()
	defer m.listenersMu.RUnlock()

	for _, listener := range m.listeners {
		select {
		case listener <- e:
		default:
			// TODO: handle the case where the channel is full
			// Optional: log dropped events or block
			// log.Println("EventManager: dropping event", e.Type())
		}
	}
}

func (l EventListener) Next(ctx context.Context) (Event, error) {
	select {
	case e := <-l:
		return e, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
