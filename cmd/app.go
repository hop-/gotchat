package cmd

import (
	"fmt"
	"log"

	"github.com/hop-/gotchat/internal/app"
	"github.com/hop-/gotchat/internal/config"
	"github.com/hop-/gotchat/internal/services"
	"github.com/hop-/gotchat/internal/storage"
	"github.com/hop-/gotchat/internal/ui/tui"
	"github.com/spf13/cobra"
)

var appCmd = &cobra.Command{
	Use:   "app",
	Short: "Run the application (same as default)",
	Run: func(cmd *cobra.Command, args []string) {
		executeApp()
	},
}

func init() {
	// Flags for app command
	appCmd.Flags().IntVarP(
		&generalServerPort,
		"port", "p",
		config.GetServerPort(),
		"port on which connection listener will be started",
	)
	appCmd.Flags().StringVarP(
		&generalDataStorageFile,
		"storage",
		"s",
		config.GetDataStorageFilePath(),
		"file to store chat data and configurations",
	)
}

func executeApp() {
	application := buildApplication()

	err := application.Init()
	if err != nil {
		log.Fatal("error:", err)
	}

	application.Run()
}

func buildApplication() *app.App {

	// Create a new application builder
	builder := app.NewBuilder().
		WithEventManager(100)
	// Get the event manager from the builder
	em := builder.GetEventManager()

	// Create a new storage and set it in the builder
	storage := storage.NewStorage("file:" + generalDataStorageFile)
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

	// Create a new server
	portStr := fmt.Sprintf(":%d", generalServerPort)
	server := services.NewServer(portStr)

	// Create a new connection manager and set it in the builder
	connectionManager := services.NewConnectionManager(
		em,
		server,
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
