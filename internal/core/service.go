package core

import (
	"context"
	"sync"
)

type Service interface {
	Init() error
	Run(ctx context.Context, wg *sync.WaitGroup)
	MapEventToCommands(event Event) []Command
	Close() error
	Name() string
}
