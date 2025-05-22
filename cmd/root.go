package cmd

import (
	"log"

	"github.com/hop-/gotchat/internal/app"
	"github.com/hop-/gotchat/internal/services"
	"github.com/hop-/gotchat/internal/storage"
	"github.com/hop-/gotchat/internal/ui/tui"
)

func Execute() {
	application := buildApplication()

	err := application.Init()
	if err != nil {
		log.Fatal("error:", err)
	}

	defer application.Close()
	application.Run()
}

func buildApplication() *app.App {

	builder := app.NewBuilder().
		WithEventManager(100)
	em := builder.GetEventManager()
	ui := tui.New(em)
	builder.WithUI(ui)
	server := services.NewServer(":7665", em)
	builder.WithService(server)

	storage := storage.NewStorage("file:chat.db")

	builder.WithService(storage)

	return builder.Build()
}
