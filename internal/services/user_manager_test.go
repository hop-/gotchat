package services

import (
	"context"
	"sync"
	"testing"

	"github.com/hop-/gotchat/internal/core"
	"github.com/stretchr/testify/mock"
)

func TestNewUserManager(t *testing.T) {
	eventEmitter := core.NewMockEventEmitter(t)
	userRepo := core.NewMockRepository[core.User](t)

	um := NewUserManager(eventEmitter, userRepo)

	if um == nil {
		t.Error("Expected UserManager to be created")
		return
	}
	if um.eventEmitter != eventEmitter {
		t.Error("Expected eventEmitter to be set")
	}
	if um.userRepo != userRepo {
		t.Error("Expected userRepo to be set")
	}
}

func TestUserManager_Name(t *testing.T) {
	um := &UserManager{}
	if um.Name() != "UserManager" {
		t.Errorf("Expected Name() to return 'UserManager', got %s", um.Name())
	}
}

func TestUserManager_Init(t *testing.T) {
	um := &UserManager{}
	if err := um.Init(); err != nil {
		t.Errorf("Expected Init() to return nil, got %v", err)
	}
}

func TestUserManager_Close(t *testing.T) {
	um := &UserManager{}
	if err := um.Close(); err != nil {
		t.Errorf("Expected Close() to return nil, got %v", err)
	}
}

func TestUserManager_Run(t *testing.T) {
	um := &UserManager{}
	ctx := context.Background()
	var wg sync.WaitGroup

	// This should not panic or block
	um.Run(ctx, &wg)
}

func TestUserManager_CreateUser(t *testing.T) {
	eventEmitter := core.NewMockEventEmitter(t)
	userRepo := core.NewMockRepository[core.User](t)

	// Setup mock expectations
	userRepo.On("Create", mock.AnythingOfType("*core.User")).Return(nil).Run(func(args mock.Arguments) {
		user := args[0].(*core.User)
		user.Id = 1 // Simulate setting ID
	})
	eventEmitter.On("Emit", mock.Anything).Return()

	um := NewUserManager(eventEmitter, userRepo)

	createdUser, err := um.CreateUser("testuser", "password123")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if createdUser == nil {
		t.Error("Expected user to be created")
		return
	}
	if createdUser.Name != "testuser" {
		t.Errorf("Expected user name to be 'testuser', got %s", createdUser.Name)
	}
}

func TestUserManager_CreateUser_InvalidInput(t *testing.T) {
	eventEmitter := core.NewMockEventEmitter(t)
	userRepo := core.NewMockRepository[core.User](t)

	um := NewUserManager(eventEmitter, userRepo)

	tests := []struct {
		name     string
		username string
		password string
	}{
		{"empty username", "", "password"},
		{"empty password", "username", ""},
		{"both empty", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := um.CreateUser(tt.username, tt.password)
			if err != ErrorInvalidInput {
				t.Errorf("Expected ErrorInvalidInput, got %v", err)
			}
			if user != nil {
				t.Error("Expected user to be nil")
			}
		})
	}
}

func TestUserManager_GetUserById(t *testing.T) {
	eventEmitter := core.NewMockEventEmitter(t)
	userRepo := core.NewMockRepository[core.User](t)

	testUser := &core.User{
		BaseEntity: core.BaseEntity{Id: 1},
		Name:       "testuser",
	}

	// First setup create user mock
	userRepo.On("Create", mock.AnythingOfType("*core.User")).Return(nil).Run(func(args mock.Arguments) {
		user := args[0].(*core.User)
		user.Id = 1
	})
	eventEmitter.On("Emit", mock.Anything).Return()

	// Setup GetOne mock
	userRepo.On("GetOne", 1).Return(testUser, nil)

	um := NewUserManager(eventEmitter, userRepo)

	// Create a user first
	createdUser, _ := um.CreateUser("testuser", "password123")

	user, err := um.GetUserById(createdUser.Id)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if user == nil {
		t.Error("Expected user to be found")
		return
	}
	if user.Name != "testuser" {
		t.Errorf("Expected user name to be 'testuser', got %s", user.Name)
	}
}

func TestUserManager_GetUserById_NotFound(t *testing.T) {
	eventEmitter := core.NewMockEventEmitter(t)
	userRepo := core.NewMockRepository[core.User](t)

	userRepo.On("GetOne", 999).Return((*core.User)(nil), nil)

	um := NewUserManager(eventEmitter, userRepo)

	user, err := um.GetUserById(999)
	if err != ErrNotFound {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
	if user != nil {
		t.Error("Expected user to be nil")
	}
}

func TestUserManager_GetUserByUniqueId(t *testing.T) {
	eventEmitter := core.NewMockEventEmitter(t)
	userRepo := core.NewMockRepository[core.User](t)

	testUser := &core.User{
		BaseEntity: core.BaseEntity{Id: 1},
		Name:       "testuser",
		UniqueId:   "unique123",
	}

	// Setup create user mock
	userRepo.On("Create", mock.AnythingOfType("*core.User")).Return(nil).Run(func(args mock.Arguments) {
		user := args[0].(*core.User)
		user.Id = 1
		user.UniqueId = "unique123"
	})
	eventEmitter.On("Emit", mock.Anything).Return()

	// Setup GetOneBy mock
	userRepo.On("GetOneBy", "unique_id", "unique123").Return(testUser, nil)

	um := NewUserManager(eventEmitter, userRepo)

	// Create a user first
	createdUser, _ := um.CreateUser("testuser", "password123")

	user, err := um.GetUserByUniqueId(createdUser.UniqueId)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if user == nil {
		t.Error("Expected user to be found")
		return
	}
	if user.Name != "testuser" {
		t.Errorf("Expected user name to be 'testuser', got %s", user.Name)
	}
}

func TestUserManager_UpdateUser(t *testing.T) {
	eventEmitter := core.NewMockEventEmitter(t)
	userRepo := core.NewMockRepository[core.User](t)

	// Setup create user mock
	userRepo.On("Create", mock.AnythingOfType("*core.User")).Return(nil).Run(func(args mock.Arguments) {
		user := args[0].(*core.User)
		user.Id = 1
	})

	// Setup update user mock
	userRepo.On("Update", mock.AnythingOfType("*core.User")).Return(nil)

	eventEmitter.On("Emit", mock.Anything).Return()

	um := NewUserManager(eventEmitter, userRepo)

	// Create a user first
	user, _ := um.CreateUser("testuser", "password123")
	user.Name = "updateduser"

	err := um.UpdateUser(user)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestUserManager_UpdateUser_NilUser(t *testing.T) {
	eventEmitter := core.NewMockEventEmitter(t)
	userRepo := core.NewMockRepository[core.User](t)

	um := NewUserManager(eventEmitter, userRepo)

	err := um.UpdateUser(nil)
	if err != ErrorInvalidInput {
		t.Errorf("Expected ErrorInvalidInput, got %v", err)
	}
}

func TestUserManager_LoginUser(t *testing.T) {
	eventEmitter := core.NewMockEventEmitter(t)
	userRepo := core.NewMockRepository[core.User](t)

	// Setup create user mock
	userRepo.On("Create", mock.AnythingOfType("*core.User")).Return(nil).Run(func(args mock.Arguments) {
		user := args[0].(*core.User)
		user.Id = 1
	})

	// Setup update user mock for login
	userRepo.On("Update", mock.AnythingOfType("*core.User")).Return(nil)

	eventEmitter.On("Emit", mock.Anything).Return()

	um := NewUserManager(eventEmitter, userRepo)

	// Create a user first
	user, _ := um.CreateUser("testuser", "password123")

	loggedUser, err := um.LoginUser(user, "password123")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if loggedUser == nil {
		t.Error("Expected user to be returned")
		return
	}
	if loggedUser.LastLogin.IsZero() {
		t.Error("Expected LastLogin to be updated")
	}
}

func TestUserManager_LoginUser_InvalidPassword(t *testing.T) {
	eventEmitter := core.NewMockEventEmitter(t)
	userRepo := core.NewMockRepository[core.User](t)

	// Setup create user mock
	userRepo.On("Create", mock.AnythingOfType("*core.User")).Return(nil).Run(func(args mock.Arguments) {
		user := args[0].(*core.User)
		user.Id = 1
	})
	eventEmitter.On("Emit", mock.Anything).Return()

	um := NewUserManager(eventEmitter, userRepo)

	// Create a user first
	user, _ := um.CreateUser("testuser", "password123")

	loggedUser, err := um.LoginUser(user, "wrongpassword")
	if err != ErrorInvalidCredentials {
		t.Errorf("Expected ErrorInvalidCredentials, got %v", err)
	}
	if loggedUser != nil {
		t.Error("Expected user to be nil")
	}
}

func TestUserManager_LoginUser_NilUser(t *testing.T) {
	eventEmitter := core.NewMockEventEmitter(t)
	userRepo := core.NewMockRepository[core.User](t)

	um := NewUserManager(eventEmitter, userRepo)

	loggedUser, err := um.LoginUser(nil, "password")
	if err != ErrorInvalidInput {
		t.Errorf("Expected ErrorInvalidInput, got %v", err)
	}
	if loggedUser != nil {
		t.Error("Expected user to be nil")
	}
}

func TestUserManager_GetAllUsers(t *testing.T) {
	eventEmitter := core.NewMockEventEmitter(t)
	userRepo := core.NewMockRepository[core.User](t)

	testUsers := []*core.User{
		{BaseEntity: core.BaseEntity{Id: 1}, Name: "user1"},
		{BaseEntity: core.BaseEntity{Id: 2}, Name: "user2"},
	}

	// Setup create user mocks
	userRepo.On("Create", mock.AnythingOfType("*core.User")).Return(nil).Run(func(args mock.Arguments) {
		user := args[0].(*core.User)
		if user.Name == "user1" {
			user.Id = 1
		} else {
			user.Id = 2
		}
	})
	eventEmitter.On("Emit", mock.Anything).Return()

	// Setup GetAll mock
	userRepo.On("GetAll").Return(testUsers, nil)

	um := NewUserManager(eventEmitter, userRepo)

	// Create some users
	um.CreateUser("user1", "password123")
	um.CreateUser("user2", "password123")

	users, err := um.GetAllUsers()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(users) != 2 {
		t.Errorf("Expected 2 users, got %d", len(users))
	}
}

func TestUserManager_GetAllUsers_Empty(t *testing.T) {
	eventEmitter := core.NewMockEventEmitter(t)
	userRepo := core.NewMockRepository[core.User](t)

	userRepo.On("GetAll").Return([]*core.User{}, nil)

	um := NewUserManager(eventEmitter, userRepo)

	users, err := um.GetAllUsers()
	if err != ErrNotFound {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
	if users != nil {
		t.Error("Expected users to be nil")
	}
}
