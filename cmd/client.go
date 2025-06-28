package cmd

import (
	"log"

	"github.com/hop-/gotchat/internal/app"
	"github.com/hop-/gotchat/internal/services"
	"github.com/hop-/gotchat/internal/storage"
	"github.com/hop-/gotchat/internal/ui/tui"
	"github.com/spf13/cobra"
)

var (
	clientCmd = &cobra.Command{
		Use:   "client",
		Short: "Run client-only mode",
		Long:  `Run the app in a client-only mode that does not require a server. Useful for testing or local development.`,
		Run: func(cmd *cobra.Command, args []string) {
			executeClient()
		},
	}
)

func init() {
}

func executeClient() {
	application := buildApplicationWithoutServer()

	err := application.Init()
	if err != nil {
		log.Fatal("error:", err)
	}

	application.Run()
}

func buildApplicationWithoutServer() *app.App {

	// Create a new application builder
	builder := app.NewBuilder().
		WithEventManager(100)
	// Get the event manager from the builder
	em := builder.GetEventManager()

	// Create a new storage and set it in the builder
	storage := storage.NewStorage("file:chat.db")
	builder.WithService(storage)

	// Create a new user manager service and set it in the builder
	userManager := services.NewUserManager(em, storage.GetUserRepository())
	builder.WithService(userManager)

	// Create a new chat manager service and set it in the builder
	chatManager := services.NewChatManager(
		em,
		userManager,
		storage.GetChannelRepository(),
		storage.GetAttendanceRepository(),
		storage.GetMessageRepository(),
	)
	builder.WithService(chatManager)

	// Create a new connection manager and set it in the builder
	connectionManager := services.NewConnectionManager(
		em,
		nil, // No server for client mode
	)

	builder.WithService(connectionManager)

	// Create a new UI and set it in the builder
	ui := tui.New(
		em,
		userManager,
		chatManager,
	)
	builder.WithUI(ui)

	return builder.Build()
}
