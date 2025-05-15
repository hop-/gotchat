package core

import "fmt"

type ServiceContainer struct {
	services []Service
}

func NewContainer() *ServiceContainer {
	return &ServiceContainer{}
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

func (c *ServiceContainer) RunAll() error {
	for _, s := range c.services {
		go s.Run()
	}

	return nil
}

func (c *ServiceContainer) CloseAll() error {
	for _, s := range c.services {
		if err := s.Close(); err != nil {
			return fmt.Errorf("service %s failed to close: %w", s.Name(), err)
		}
	}

	return nil
}
