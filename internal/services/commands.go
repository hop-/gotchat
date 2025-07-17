package services

import (
	"context"

	"github.com/hop-/gotchat/internal/core"
)

type Connect struct {
	cm      *ConnectionManager
	address string
}

func (c *Connect) Execute(ctx context.Context) ([]core.Event, error) {
	c.cm.Connect(c.address)

	return nil, nil
}

type ChangeUserController struct {
	cm   *ConnectionManager
	User *core.User
}

func (c *ChangeUserController) Execute(ctx context.Context) ([]core.Event, error) {
	c.cm.changeUserController(c.User)

	return nil, nil
}

type RemoveUserController struct {
	cm *ConnectionManager
}

func (r *RemoveUserController) Execute(ctx context.Context) ([]core.Event, error) {
	r.cm.removeUserController()

	return nil, nil
}
