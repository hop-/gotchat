package cmd

import (
	"context"
	"log"

	"github.com/hop-/gotchat/internal/app"
	"github.com/hop-/gotchat/internal/services"
	"github.com/hop-/gotchat/internal/ui/tui"
)

func Execute() {
	application := buildApplication()

	err := application.Init()
	if err != nil {
		log.Fatal("error:", err)
	}

	defer application.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	application.Run(ctx)
}

func buildApplication() *app.App {
	ui := tui.New()

	builder := app.NewBuilder().
		WithEventManager(100).
		WithUI(ui)
	em := builder.GetEventManager()
	server := services.NewServer("localhost:7665", em)
	builder.WithService(server)

	return builder.Build()
}
