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
}

func (a *App) Init() error {
	return a.services.InitAll()
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

		logic.Handle(event)

		// err = a.ui.Send(update)
	}
}

func (a *App) Close() {
	// TODO: Handle errors
	a.services.CloseAll()

	// TODO: Handle error
	a.ui.Close()
}
