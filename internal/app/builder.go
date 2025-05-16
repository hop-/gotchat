package app

import (
	"github.com/hop-/gotchat/internal/core"
	"github.com/hop-/gotchat/internal/logic"
	"github.com/hop-/gotchat/internal/ui"
)

type Builder struct {
	em       *core.EventManager
	ui       ui.UI
	services []core.Service
	appLogic *logic.AppLogic
}

func NewBuilder() *Builder {
	return &Builder{
		services: make([]core.Service, 0),
		appLogic: logic.New(),
	}
}

func (b *Builder) WithEventManager(bufferSize int) *Builder {
	b.em = core.NewEventManager(bufferSize)

	return b
}

func (b *Builder) WithUI(ui ui.UI) *Builder {
	b.ui = ui

	return b
}

func (b *Builder) WithService(s core.Service) *Builder {
	b.services = append(b.services, s)

	return b
}

func (b *Builder) GetEventManager() *core.EventManager {
	return b.em
}

func (b *Builder) Build() *App {
	container := core.NewContainer()
	for _, s := range b.services {
		container.Register(s)
	}

	if b.em == nil || b.ui == nil {
		panic("EventManager and UI must be set")
	}

	return &App{
		b.em,
		container,
		b.ui,
		b.appLogic,
	}
}
