package ui

import (
	"context"
	"sync"
)

type UI interface {
	Init() error
	Run(ctx context.Context, wg *sync.WaitGroup) error
	Close() error
}
