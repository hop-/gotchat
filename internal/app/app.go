package app

import (
	"context"

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

func (a *App) Run(ctx context.Context) error {
	go a.ui.Run()
	a.services.RunAll()

	for {
		event, err := a.eventManager.Next(ctx)
		if err != nil {
			// TODO: Handle error
			continue
		}

		a.logic.Handle(event)

		// err = a.ui.Send(update)
	}
}

func (app *App) Close() {
	// TODO: Handle errors
	app.services.CloseAll()

	// TODO: Handle error
	app.ui.Close()
}
