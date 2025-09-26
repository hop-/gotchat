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
	accountRepo  core.Repository[core.Account]
}

func NewUserManager(eventEmitter core.EventEmitter, userRepo core.Repository[core.User], accountRepo core.Repository[core.Account]) *UserManager {
	return &UserManager{
		eventEmitter,
		userRepo,
		accountRepo,
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
	user, err := u.userRepo.GetOneBy("UniqueId", uniqueId)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserManager) GetUserById(id uint) (*core.User, error) {
	user, err := u.userRepo.GetOne(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserManager) GetAllUsers() ([]*core.User, error) {
	return u.userRepo.GetAll()
}

func (u *UserManager) GetAllAccountUsers() ([]*core.User, error) {
	accounts, err := u.accountRepo.GetAll()
	if err != nil {
		return nil, err
	}

	if len(accounts) == 0 {
		return []*core.User{}, nil
	}

	users := make([]*core.User, 0, len(accounts))
	for _, account := range accounts {
		user, err := u.userRepo.GetOne(account.UserId)
		if err != nil {
			return nil, err
		}
		if user != nil {
			users = append(users, user)
		}
	}

	return users, nil
}

func (u *UserManager) UpdateUser(user *core.User) error {
	if user == nil {
		return ErrorInvalidInput
	}

	updated, err := u.userRepo.Update(user)
	if err != nil {
		return err
	}

	u.eventEmitter.Emit(core.UserUpdatedEvent{
		User: updated,
	})

	return nil
}

func (u *UserManager) CreateUser(name string) (*core.User, error) {
	if name == "" {
		return nil, ErrorInvalidInput
	}

	user := core.NewUser(name)
	updated, err := u.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	u.eventEmitter.Emit(core.UserCreatedEvent{
		User: updated,
	})

	return updated, nil
}

func (u *UserManager) DeleteUser(user *core.User) error {
	if user == nil {
		return ErrorInvalidInput
	}

	// Delete associated account if exists
	account, err := u.accountRepo.GetOneBy("UserId", user.Id)
	if err != nil {
		return err
	}
	if account != nil {
		if err := u.accountRepo.Delete(account.Id); err != nil {
			return err
		}
	}

	if err := u.userRepo.Delete(user.Id); err != nil {
		return err
	}

	return nil
}

func (u *UserManager) GetAccountByUser(user *core.User) (*core.Account, error) {
	if user == nil {
		return nil, ErrorInvalidInput
	}

	account, err := u.accountRepo.GetOneBy("UserId", user.Id)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (u *UserManager) CreateUserAccount(user *core.User, password string) (*core.Account, error) {
	if password == "" {
		return nil, ErrorInvalidInput
	}

	account, err := u.GetAccountByUser(user)
	if err != nil && err != core.ErrEntityNotFound {
		return nil, err
	}
	if account != nil {
		return nil, ErrorEntityExists
	}

	passwordHash, err := core.HashPassword(password)
	if err != nil {
		return nil, err
	}

	account = core.NewAccount(user.Id, passwordHash, time.Now())
	created, err := u.accountRepo.Create(account)
	if err != nil {
		return nil, err
	}

	u.eventEmitter.Emit(core.UserAccountCreatedEvent{
		User:    user,
		Account: created,
	})

	return created, nil
}

func (u *UserManager) UpdateAccount(account *core.Account) error {
	if account == nil {
		return ErrorInvalidInput
	}

	updated, err := u.accountRepo.Update(account)
	if err != nil {
		return err
	}

	u.eventEmitter.Emit(core.UserAccountUpdatedEvent{
		Account: updated,
	})

	return nil
}

func (u *UserManager) checkPassword(account *core.Account, password string) (bool, error) {
	if account == nil {
		return false, ErrorInvalidInput
	}

	return core.CheckPasswordHash(password, account.Password), nil
}

func (u *UserManager) LoginUser(user *core.User, password string) (*core.User, error) {
	if user == nil {
		return nil, ErrorInvalidInput
	}

	account, err := u.GetAccountByUser(user)
	if err != nil {
		return nil, err
	}

	isValid, err := u.checkPassword(account, password)
	if err != nil {
		return nil, err
	}
	if !isValid {
		return nil, ErrorInvalidCredentials
	}

	account.LastLogin = time.Now()
	if err := u.UpdateAccount(account); err != nil {
		return nil, err
	}

	u.eventEmitter.Emit(core.UserLoggedInEvent{
		User:    user,
		Account: account,
	})

	return user, nil
}
