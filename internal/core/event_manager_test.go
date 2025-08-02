package core

import (
	"context"
	"testing"
	"time"
)

type mockEvent struct{}

func TestEventManager_RegisterAndUnregister(t *testing.T) {
	em := NewEventManager(10)
	ctx, cancel := context.WithCancel(context.Background())

	listener := em.Register(ctx)
	if listener == nil {
		t.Fatal("Expected listener to be registered, got nil")
	}

	em.Unregister(listener)
	select {
	case _, ok := <-listener:
		if ok {
			t.Fatal("Expected listener channel to be closed")
		}
	default:
		t.Fatal("Expected listener channel to be closed")
	}

	cancel()
}

func TestEventManager_Emit(t *testing.T) {
	em := NewEventManager(10)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	listener := em.Register(ctx)
	event := mockEvent{}

	em.Emit(event)

	select {
	case receivedEvent := <-listener:
		if receivedEvent != event {
			t.Fatalf("Expected event %v, got %v", event, receivedEvent)
		}
	case <-time.After(time.Second):
		t.Fatal("Expected event to be received, but timed out")
	}
}

func TestEventListener_Next(t *testing.T) {
	listener := make(EventListener, 1)
	event := mockEvent{}
	listener <- event

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	receivedEvent, err := listener.Next(ctx)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if receivedEvent != event {
		t.Fatalf("Expected event %v, got %v", event, receivedEvent)
	}
}

func TestEventListener_Next_ContextCancelled(t *testing.T) {
	listener := make(EventListener, 1)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := listener.Next(ctx)
	if err == nil {
		t.Fatal("Expected error due to cancelled context, got nil")
	}
}
