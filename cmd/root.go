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

	// Create a new server and set it in the builder
	server := services.NewServer(":7665", em)
	builder.WithService(server)

	// Create a new UI and set it in the builder
	ui := tui.New(
		em,
		userManager,
		chatManager,
	)
	builder.WithUi(ui)

	return builder.Build()
}
