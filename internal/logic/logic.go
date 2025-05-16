package logic

import "github.com/hop-/gotchat/internal/core"

type AppLogic struct {
}

func New() *AppLogic {
	return &AppLogic{}
}

func (l *AppLogic) Handle(e core.Event) {
	// TODO
}

// TODO: get data from the app logic
// func (l *AppLogic) GetSomeData() SomeData { return l.someData }
