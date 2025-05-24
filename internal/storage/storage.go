package storage

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	"github.com/hop-/gotchat/internal/core"
	_ "modernc.org/sqlite" // SQLite driver
)

type StorageDb interface {
	Db() *sql.DB
}

type Storage struct {
	path string
	db   *sql.DB
}

func NewStorage(path string) *Storage {
	return &Storage{path, nil}
}

func (s *Storage) Db() *sql.DB {
	return s.db
}

func (s *Storage) Init() error {
	// Start the server
	if s.db != nil {
		return fmt.Errorf("server is already running")
	}

	db, err := sql.Open("sqlite", s.path)
	if err != nil {
		return err
	}

	s.db = db

	return s.createTables()
}

func (s *Storage) Run(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	<-ctx.Done()
}

func (s *Storage) Close() error {
	if s.db == nil {
		return nil
	}

	err := s.db.Close()
	if err != nil {
		return err
	}

	s.db = nil
	return nil
}

func (s *Storage) GetUserRepository() core.Repository[core.User] {
	return newUserRepository(s)
}

func (s *Storage) GetChannelRepository() core.Repository[core.Channel] {
	return newChannelRepository(s)
}

func (s *Storage) Name() string {
	return "Storage"
}

func (s *Storage) createTables() error {
	err := createUserTable(s.db)

	if err != nil {
		return err
	}

	err = createChannelTable(s.db)
	if err != nil {
		return err
	}

	return nil
}

func createUserTable(db *sql.DB) error {
	// Create the users table if it doesn't exist
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		unique_id TEXT UNIQUE,
		name TEXT,
		last_login DATETIME
	)`)

	return err
}

func createChannelTable(db *sql.DB) error {
	// Create the channels table if it doesn't exist
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS channels (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		unique_id TEXT UNIQUE,
		name TEXT
	)`)

	return err
}

func isFieldExist[T core.Entity](field string) bool {
	fields := core.GetFieldNamesOfEntity[T]()

	for _, filedName := range fields {
		if filedName == field {
			return true
		}
	}

	return false
}
