package storage

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

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

	return nil
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

func (s *Storage) Name() string {
	return "Storage"
}
