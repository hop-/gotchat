package storage

import (
	"context"
	"fmt"
	"sync"

	"github.com/hop-/gotchat/internal/core"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type StorageDb interface {
	Db() *gorm.DB
}

type Storage struct {
	path string
	db   *gorm.DB

	// Repositories
	accountRepository GormRepo[Account, core.Account]
	userRepo          GormRepo[User, core.User]
	channelRepo       GormRepo[Channel, core.Channel]
	attendanceRepo    GormRepo[Attendance, core.Attendance]
	messageRepo       GormRepo[Message, core.Message]
	connectionRepo    GormRepo[ConnectionDetails, core.ConnectionDetails]
}

func NewStorage(path string) *Storage {
	return &Storage{path: path}
}

func (s *Storage) Db() *gorm.DB {
	return s.db
}

func (s *Storage) Init() error {
	// Start the service
	if s.db != nil {
		return fmt.Errorf("server is already running")
	}

	db, err := gorm.Open(sqlite.Open(s.path), &gorm.Config{
		// Silent logger
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return err
	}

	s.db = db

	err = s.configureDatabase()
	if err != nil {
		s.db = nil

		return err
	}

	s.db.AutoMigrate(&User{}, &Account{}, &ConnectionDetails{}, &Channel{}, &Attendance{}, &Message{})

	// Initialize repositories
	err = s.initRepositories()
	if err != nil {
		s.db = nil

		return err
	}

	return nil
}

func (s *Storage) Run(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	<-ctx.Done()
}

func (s *Storage) MapEventToCommands(event core.Event) []core.Command {
	// This method is not used in the Storage service, so we return an empty slice.
	return nil
}

func (s *Storage) Close() error {
	if s.db == nil {
		return nil
	}
	s.db = nil

	return nil
}

func (s *Storage) GetAccountRepository() core.Repository[core.Account] {
	if s.accountRepository == nil {
		s.accountRepository = newAccountRepository(s)
	}

	return s.accountRepository
}

func (s *Storage) GetUserRepository() core.Repository[core.User] {
	if s.userRepo == nil {
		s.userRepo = newUserRepository(s)
	}

	return s.userRepo
}

func (s *Storage) GetConnectionDetailsRepository() core.Repository[core.ConnectionDetails] {
	if s.connectionRepo == nil {
		s.connectionRepo = newConnectionDetailsRepository(s)
	}

	return s.connectionRepo
}

func (s *Storage) GetChannelRepository() core.Repository[core.Channel] {
	if s.channelRepo == nil {
		s.channelRepo = newChannelRepository(s)
	}

	return s.channelRepo
}

func (s *Storage) GetAttendanceRepository() core.Repository[core.Attendance] {
	if s.attendanceRepo == nil {
		s.attendanceRepo = newAttendanceRepository(s)
	}

	return s.attendanceRepo
}

func (s *Storage) GetMessageRepository() core.Repository[core.Message] {
	if s.messageRepo == nil {
		s.messageRepo = newMessageRepository(s)
	}

	return s.messageRepo
}

func (s *Storage) Name() string {
	return "Storage"
}

func (s *Storage) configureDatabase() error {
	db, err := s.db.DB()
	if err != nil {
		return err
	}

	// Enable foreign key constraints
	_, err = db.Exec(`PRAGMA foreign_keys = ON;`)
	if err != nil {
		return err
	}

	// Set journal mode to WAL for better concurrency
	_, err = db.Exec(`PRAGMA journal_mode = WAL;`)
	if err != nil {
		return err
	}

	// Configure busy timeout to handle database locks
	_, err = db.Exec(`PRAGMA busy_timeout = 5000;`) // 5000 milliseconds
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) initRepositories() error {
	// Initialize all repositories
	s.GetUserRepository()
	s.GetAccountRepository()
	s.GetConnectionDetailsRepository()
	s.GetChannelRepository()
	s.GetAttendanceRepository()
	s.GetMessageRepository()

	// Call Init on each repository and handle errors
	err := s.userRepo.Init()
	if err != nil {
		return err
	}

	err = s.accountRepository.Init()
	if err != nil {
		return err
	}

	err = s.connectionRepo.Init()
	if err != nil {
		return err
	}

	err = s.channelRepo.Init()
	if err != nil {
		return err
	}

	err = s.attendanceRepo.Init()
	if err != nil {
		return err
	}

	err = s.messageRepo.Init()
	if err != nil {
		return err
	}

	return nil
}
