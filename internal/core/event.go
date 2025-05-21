package core

type EventType string

type Event interface{}

type BaseEvent struct {
	EventType EventType
	Data      any
}

func (e BaseEvent) Type() EventType {
	return e.EventType
}

func (e BaseEvent) Payload() any {
	return e.Data
}
