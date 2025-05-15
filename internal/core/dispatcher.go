package core

import "sync"

type Handler func(Event)

type Dispatcher struct {
	mu       sync.RWMutex
	handlers map[EventType][]Handler
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{handlers: make(map[EventType][]Handler)}
}

func (d *Dispatcher) Subscribe(eventType EventType, handler Handler) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.handlers[eventType] = append(d.handlers[eventType], handler)
}

func (d *Dispatcher) Emit(e Event) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	for _, h := range d.handlers[e.Type()] {
		go h(e)
	}
}
