package app

import (
	"context"
	"testing"

	"github.com/hop-/gotchat/internal/core"
	"github.com/hop-/gotchat/internal/ui"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestApp_Init(t *testing.T) {
	mockServices := core.NewMockServiceDispatcher(t)
	mockServices.On("InitAll").Return(nil)

	app := &App{
		services: mockServices,
	}

	err := app.Init()
	assert.NoError(t, err)
	mockServices.AssertCalled(t, "InitAll")
}

func TestApp_mapEventToCommands(t *testing.T) {
	mockServices := core.NewMockServiceDispatcher(t)
	mockService := core.NewMockService(t)

	mockServices.On("GetAll").Return([]core.Service{mockService})
	mockService.On("MapEventToCommands", mock.Anything).Return([]core.Command{})

	app := &App{
		services: mockServices,
	}

	commands := app.mapEventToCommands(core.QuitEvent{})
	assert.NotNil(t, commands)
	mockServices.AssertCalled(t, "GetAll")
	mockService.AssertCalled(t, "MapEventToCommands", mock.Anything)
}

func TestApp_executeCommands(t *testing.T) {
	mockEventManager := core.NewMockEventDispatcher(t)
	mockCommand := core.NewMockCommand(t)

	mockCommand.On("Execute", mock.Anything).Return([]core.Event{}, nil)

	app := &App{
		eventManager: mockEventManager,
	}

	app.executeCommands([]core.Command{mockCommand}, context.Background())

	mockCommand.AssertCalled(t, "Execute", mock.Anything)
}

func TestApp_close(t *testing.T) {
	mockServices := core.NewMockServiceDispatcher(t)
	mockUI := ui.NewMockUI(t)

	mockServices.On("CloseAll").Return([]error{})
	mockUI.On("Close").Return(nil)

	app := &App{
		services: mockServices,
		ui:       mockUI,
	}

	app.close()

	mockServices.AssertCalled(t, "CloseAll")
	mockUI.AssertCalled(t, "Close")
}
