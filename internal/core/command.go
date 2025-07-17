package core

import (
	"context"
)

type Command interface {
	Execute(ctx context.Context) ([]Event, error)
}
