package app

import (
	"context"
	"sync"

	"github.com/hop-/gotchat/internal/core"
	"github.com/hop-/gotchat/internal/log"
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

func (a *App) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	eventListner := a.eventManager.Register(ctx)

	wg := sync.WaitGroup{}

	// Run all services and UI in separate goroutines
	a.services.RunAll(ctx, &wg)
	go a.ui.Run(ctx, &wg)

	// Run the event loop
	isRunning := true
	for isRunning {
		event, err := eventListner.Next(ctx)
		if err != nil {
			// TODO: Handle error if needed
			log.Errorf("Failed to get next event: %v\n", err)
			continue
		}

		commands := a.mapEventToCommands(event)

		a.executeCommands(commands, ctx)

		switch event.(type) {
		case core.QuitEvent:
			isRunning = false
		}
	}

	//cancel all services and UI
	cancel()

	// Close the application gracefully
	a.close()

	// Wait for all goroutines to finish
	wg.Wait()
}

func (a *App) mapEventToCommands(event core.Event) []core.Command {
	// TODO: Implement priority handling if needed
	commands := make([]core.Command, 0)
	for _, service := range a.services.GetAll() {
		cmds := service.MapEventToCommands(event)
		commands = append(commands, cmds...)
	}

	return commands
}

func (a *App) executeCommands(commands []core.Command, ctx context.Context) {
	var events []core.Event
	for _, cmd := range commands {
		cmdEvents, err := cmd.Execute(ctx)
		if err != nil {
			log.Errorf("Failed to execute command: %v\n", err)

			continue
		}
		events = append(events, cmdEvents...)
	}

	for _, e := range events {
		a.eventManager.Emit(e)
	}
}

func (a *App) close() {
	errs := a.services.CloseAll()
	if len(errs) > 0 {
		for _, err := range errs {
			log.Errorf("Failed to close service: %v\n", err)
		}
	}

	err := a.ui.Close()
	if err != nil {
		log.Errorf("Failed to close UI: %v\n", err)
	}
}
