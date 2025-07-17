package services

import (
	"context"
	"sync"
	"time"

	"github.com/hop-/gotchat/internal/core"
)

type UserManager struct {
	eventEmitter core.EventEmitter
	userRepo     core.Repository[core.User]
}

func NewUserManager(eventEmitter core.EventEmitter, userRepo core.Repository[core.User]) *UserManager {
	return &UserManager{
		eventEmitter,
		userRepo,
	}
}

func (u *UserManager) Init() error {
	return nil
}

func (u *UserManager) Name() string {
	return "UserManager"
}

func (u *UserManager) Run(ctx context.Context, wg *sync.WaitGroup) {
	// This service does not run any background tasks.
}

func (u *UserManager) MapEventToCommands(event core.Event) []core.Command {
	// TODO
	return nil
}

func (u *UserManager) Close() error {
	return nil
}

func (u *UserManager) GetUserByUniqueId(uniqueId string) (*core.User, error) {
	user, err := u.userRepo.GetOneBy("unique_id", uniqueId)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrNotFound
	}

	return user, nil
}

func (u *UserManager) GetUserById(id int) (*core.User, error) {
	user, err := u.userRepo.GetOne(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrNotFound
	}

	return user, nil
}

func (u *UserManager) GetAllUsers() ([]*core.User, error) {
	users, err := u.userRepo.GetAll()
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, ErrNotFound
	}

	return users, nil
}

func (u *UserManager) UpdateUser(user *core.User) error {
	if user == nil {
		return ErrorInvalidInput
	}

	if err := u.userRepo.Update(user); err != nil {
		return err
	}

	u.eventEmitter.Emit(core.UserUpdatedEvent{
		User: user,
	})

	return nil
}

func (u *UserManager) CreateUser(name string, password string) (*core.User, error) {
	if name == "" || password == "" {
		return nil, ErrorInvalidInput
	}

	passwordHash, err := core.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := core.NewUser(name, passwordHash)
	if err := u.userRepo.Create(user); err != nil {
		return nil, err
	}

	u.eventEmitter.Emit(core.UserCreatedEvent{
		User: user,
	})

	return user, nil
}

func (u *UserManager) checkPasswordByUser(user *core.User, password string) bool {
	if user == nil {
		return false
	}

	return core.CheckPasswordHash(password, user.Password)
}

func (u *UserManager) LoginUser(user *core.User, password string) (*core.User, error) {
	if user == nil {
		return nil, ErrorInvalidInput
	}

	if !u.checkPasswordByUser(user, password) {
		return nil, ErrorInvalidCredentials
	}

	user.LastLogin = time.Now()
	if err := u.UpdateUser(user); err != nil {
		return nil, err
	}

	u.eventEmitter.Emit(core.UserLoggedInEvent{
		User: user,
	})

	return user, nil
}
