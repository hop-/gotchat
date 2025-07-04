package app

import (
	"context"
	"log"
	"sync"

	"github.com/hop-/gotchat/internal/core"
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
			log.Printf("failed to get next event: %v\n", err)
			continue
		}

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

func (a *App) close() {
	errs := a.services.CloseAll()
	if len(errs) > 0 {
		for _, err := range errs {
			log.Printf("failed to close service: %v\n", err)
		}
	}

	err := a.ui.Close()
	if err != nil {
		log.Printf("failed to close UI: %v\n", err)
	}
}
