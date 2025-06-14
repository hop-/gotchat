package services

import (
	"context"
	"sync"

	"github.com/hop-/gotchat/internal/core"
)

type UserManager struct {
	em       *core.EventManager
	userRepo core.Repository[core.User]
}

func NewUserManager(em *core.EventManager, userRepo core.Repository[core.User]) *UserManager {
	return &UserManager{
		em,
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

	// TODO: Emit an event for user update

	return nil
}

func (u *UserManager) CreateUser(user *core.User) error {
	if user == nil {
		return ErrorInvalidInput
	}

	if user.UniqueId == "" {
		return ErrorInvalidInput
	}

	if user.Password == "" {
		return ErrorInvalidInput
	}

	if err := u.userRepo.Create(user); err != nil {
		return err
	}

	// TODO: Emit an event for user creation
	return nil
}

func (u *UserManager) CheckPasswordByUser(user *core.User, password string) bool {
	if user == nil {
		return false
	}

	return core.CheckPasswordHash(password, user.Password)
}
