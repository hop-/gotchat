package cmd

import (
	"github.com/hop-/gotchat/internal/app"
	"github.com/hop-/gotchat/internal/config"
	"github.com/hop-/gotchat/internal/log"
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
	// Flags for client command
	clientCmd.Flags().StringVarP(
		&generalDataStorageFile,
		"storage",
		"s",
		config.GetDataStorageFilePath(),
		"file to store chat data and configurations",
	)
}

func executeClient() {
	application := buildApplicationWithoutServer()

	err := application.Init()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	application.Run()
}

func buildApplicationWithoutServer() *app.App {

	// Create a new application builder
	builder := app.NewBuilder().
		WithEventDispatcher(100)
	// Get the event manager from the builder
	em := builder.GetEventManager()

	// Create a new storage and set it in the builder
	storage := storage.NewStorage("file:" + generalDataStorageFile)
	builder.WithService(storage)

	// Create a new user manager service and set it in the builder
	userManager := services.NewUserManager(em, storage.GetUserRepository())
	builder.WithService(userManager)

	// Create a new connection details manager and set it in the builder
	connectionDetailsManager := services.NewConnectionDetailsManager(em, storage.GetConnectionDetailsRepository())
	builder.WithService(connectionDetailsManager)

	// Create a new chat manager service and set it in the builder
	chatManager := services.NewChatManager(
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
		userManager,
		connectionDetailsManager,
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
