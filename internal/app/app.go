package app

import (
	"context"
	"sync"

	"github.com/hop-/gotchat/internal/core"
	"github.com/hop-/gotchat/internal/logic"
	"github.com/hop-/gotchat/internal/ui"
)

type App struct {
	eventManager *core.EventManager
	services     *core.ServiceContainer
	ui           ui.UI
	logic        *logic.AppLogic
}

func (app *App) Init() error {
	return app.services.InitAll()
}

func (a *App) Run() {
	ctx, cancel := context.WithCancel(context.Background())

	wg := sync.WaitGroup{}

	// Run all services and UI in separate goroutines
	a.services.RunAll(ctx, &wg)
	go a.ui.Run(ctx, &wg)

	// Run the event loop
	isRunning := true
	for isRunning {
		event, err := a.eventManager.Next(ctx)
		if err != nil {
			// TODO: Handle error
			continue
		}

		switch event.(type) {
		case core.QuitEvent:
			isRunning = false
		}

		a.logic.Handle(event)
	}

	//cancel all services and UI
	cancel()
	// Wait for all goroutines to finish
	wg.Wait()
}

func (app *App) Close() {
	// TODO: Handle errors
	app.services.CloseAll()

	// TODO: Handle error
	app.ui.Close()
}
