package core

import (
	"context"
	"sync"
	"testing"
	"time"

	mock "github.com/stretchr/testify/mock"
)

func TestServiceContainer_Register(t *testing.T) {
	container := NewContainer()

	mockService := NewMockService(t)
	mockService.On("Name").Return("MockService")

	container.Register(mockService)

	services := container.GetAll()
	if len(services) != 1 {
		t.Fatalf("expected 1 service, got %d", len(services))
	}

	if services[0].Name() != "MockService" {
		t.Fatalf("expected service name 'MockService', got '%s'", services[0].Name())
	}
}

func TestServiceContainer_GetAll(t *testing.T) {
	container := NewContainer()

	mockService1 := NewMockService(t)
	mockService1.On("Name").Return("MockService1")
	mockService2 := NewMockService(t)
	mockService2.On("Name").Return("MockService2")

	container.Register(mockService1)
	container.Register(mockService2)

	services := container.GetAll()
	if len(services) != 2 {
		t.Fatalf("expected 2 services, got %d", len(services))
	}

	if services[0].Name() != "MockService1" || services[1].Name() != "MockService2" {
		t.Fatalf("expected service names 'MockService1' and 'MockService2', got '%s' and '%s'", services[0].Name(), services[1].Name())
	}
}

func TestServiceContainer_InitAll(t *testing.T) {
	container := NewContainer()

	mockService1 := NewMockService(t)
	mockService1.On("Init").Return(nil)

	container.Register(mockService1)

	mockService2 := NewMockService(t)
	mockService2.On("Init").Return(nil)

	container.Register(mockService2)

	if err := container.InitAll(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	mockService1.AssertCalled(t, "Init")
	mockService2.AssertCalled(t, "Init")
}

func TestServiceContainer_RunAll(t *testing.T) {
	container := NewContainer()

	mockService := NewMockService(t)
	mockService.On("Run", mock.Anything, mock.Anything).Return(nil)

	container.Register(mockService)

	wg := &sync.WaitGroup{}
	ctx := context.Background()

	if err := container.RunAll(ctx, wg); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	time.Sleep(100 * time.Millisecond) // Allow goroutines to start
	mockService.AssertCalled(t, "Run", ctx, wg)
	wg.Wait()
}

func TestServiceContainer_CloseAll(t *testing.T) {
	container := NewContainer()

	mockService1 := NewMockService(t)
	mockService1.On("Close").Return(nil)

	container.Register(mockService1)

	mockService2 := NewMockService(t)
	mockService2.On("Close").Return(nil)

	container.Register(mockService2)

	errs := container.CloseAll()
	if len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}

	mockService1.AssertCalled(t, "Close")
	mockService2.AssertCalled(t, "Close")
}
