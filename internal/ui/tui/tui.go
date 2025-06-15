package tui

import (
	"context"
	"fmt"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hop-/gotchat/internal/core"
	"github.com/hop-/gotchat/internal/services"
)

type Tui struct {
	p  *tea.Program
	em *core.EventManager
}

func New(
	em *core.EventManager,
	userManager *services.UserManager,
	chatManager *services.ChatManager,
) *Tui {
	rootModel := newRootModel(newUsersListModel(userManager, chatManager), em)
	p := tea.NewProgram(rootModel, tea.WithAltScreen())

	return &Tui{p, em}
}

func (ui *Tui) Init() error {
	// Not needed for Bubble Tea
	return nil
}

func (ui *Tui) Run(ctx context.Context, wg *sync.WaitGroup) error {
	wg.Add(1)
	defer wg.Done()

	externalQuit := false
	// Quit the program when the context is done
	go func() {
		<-ctx.Done()
		externalQuit = true
		ui.p.Quit()
	}()

	listener := ui.em.Register(ctx)

	// Filter events and send to TUI program
	go ui.runEventFilter(listener)

	// Run the TUI program
	_, err := ui.p.Run()

	// Send a quit event to the event manager when tui is closed from inside
	if !externalQuit {
		ui.em.Emit(core.QuitEvent{})
	}

	return err
}

func (ui *Tui) Close() error {
	ui.p.Quit()

	return nil
}

func (ui *Tui) runEventFilter(listener core.EventListener) {
	for event := range listener {
		switch event := event.(type) {
		case core.NewMessageEvent:
			// TODO
			fmt.Println("New message event received:", event.Message)
		}
	}
}
