package core

type Service interface {
	Init() error
	Run()
	Close() error
	Name() string
}
