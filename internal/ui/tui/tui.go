package tui

import (
	"context"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hop-/gotchat/internal/core"
	"github.com/hop-/gotchat/internal/services"
)

type Tui struct {
	p       *tea.Program
	emitter core.EventEmitter
}

func New(
	em core.EventEmitter,
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

	_, err := ui.p.Run()

	// Send a quit event to the event manager when tui is closed from inside
	if !externalQuit {
		ui.emitter.Emit(core.QuitEvent{})
	}

	return err
}

func (ui *Tui) Close() error {
	ui.p.Quit()

	return nil
}
