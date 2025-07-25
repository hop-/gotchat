package core

import (
	"context"
	"fmt"
	"sync"
)

type ServiceContainer struct {
	services []Service
}

func NewContainer() *ServiceContainer {
	return &ServiceContainer{}
}

func (c *ServiceContainer) GetAll() []Service {
	return c.services
}

func (c *ServiceContainer) Register(s Service) {
	c.services = append(c.services, s)
}

func (c *ServiceContainer) InitAll() error {
	for _, s := range c.services {
		if err := s.Init(); err != nil {
			return fmt.Errorf("service %s failed to init: %w", s.Name(), err)
		}
	}

	return nil
}

func (c *ServiceContainer) RunAll(ctx context.Context, wg *sync.WaitGroup) error {
	for _, s := range c.services {
		go s.Run(ctx, wg)
	}

	return nil
}

func (c *ServiceContainer) CloseAll() []error {
	errs := make([]error, 0, len(c.services))
	for _, s := range c.services {
		if err := s.Close(); err != nil {
			errs = append(errs, fmt.Errorf("service %s failed to close: %w", s.Name(), err))
		}
	}

	if len(errs) != 0 {
		return errs
	}

	return nil
}
