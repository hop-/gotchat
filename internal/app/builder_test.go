package app

import (
	"testing"

	"github.com/hop-/gotchat/internal/core"
	"github.com/hop-/gotchat/internal/log"
	"github.com/hop-/gotchat/internal/ui"
)

func TestBuilder_WithEventDispatcher(t *testing.T) {
	builder := NewBuilder()
	bufferSize := 10
	builder.WithEventDispatcher(bufferSize)

	if builder.GetEventManager() == nil {
		t.Error("Expected EventDispatcher to be initialized, got nil")
	}
}

func TestBuilder_WithUI(t *testing.T) {
	builder := NewBuilder()
	mockUI := ui.NewMockUI(t)
	builder.WithUI(mockUI)

	if builder.Build().ui == nil {
		t.Error("Expected UI to be set, got nil")
	}
}

func TestBuilder_WithService(t *testing.T) {
	builder := NewBuilder()
	mockService := core.NewMockService(t)
	builder.WithService(mockService)

	appInstance := builder.Build()
	if len(appInstance.services.GetAll()) != 1 {
		t.Errorf("Expected 1 service to be registered, got %d", len(appInstance.services.GetAll()))
	}
}

func TestBuilder_Build(t *testing.T) {
	builder := NewBuilder()
	mockUI := ui.NewMockUI(t)
	mockService := core.NewMockService(t)

	builder.WithEventDispatcher(10).WithUI(mockUI).WithService(mockService)
	appInstance := builder.Build()

	if appInstance.eventManager == nil {
		t.Error("Expected EventManager to be set, got nil")
	}
	if appInstance.ui == nil {
		t.Error("Expected UI to be set, got nil")
	}
	if len(appInstance.services.GetAll()) != 1 {
		t.Errorf("Expected 1 service to be registered, got %d", len(appInstance.services.GetAll()))
	}
}

func TestBuilder_Build_MissingDependencies(t *testing.T) {
	log.Configure().
		Level(log.FATAL).
		Init()
	defer log.Close()

	builder := NewBuilder()

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic due to missing dependencies, but no panic occurred")
		}
	}()

	builder.Build()
}
