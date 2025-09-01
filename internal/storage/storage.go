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

	// Repositories
	userRepo       core.Repository[core.User]
	channelRepo    core.Repository[core.Channel]
	attendanceRepo core.Repository[core.Attendance]
	messageRepo    core.Repository[core.Message]
	connectionRepo core.Repository[core.ConnectionDetails]
}

func NewStorage(path string) *Storage {
	return &Storage{path, nil, nil, nil, nil, nil, nil}
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

	err = s.configureDatabase()
	if err != nil {
		s.db.Close()
		s.db = nil
		return err
	}

	return s.createTables()
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

	err := s.db.Close()
	if err != nil {
		return err
	}

	s.db = nil
	return nil
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
	// Enable foreign key constraints
	_, err := s.db.Exec(`PRAGMA foreign_keys = ON;`)
	if err != nil {
		return err
	}

	// Set journal mode to WAL for better concurrency
	_, err = s.db.Exec(`PRAGMA journal_mode = WAL;`)
	if err != nil {
		return err
	}

	// Configure busy timeout to handle database locks
	_, err = s.db.Exec(`PRAGMA busy_timeout = 5000;`) // 5000 milliseconds
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) createTables() error {
	err := createUserTable(s.db)
	if err != nil {
		return err
	}

	err = createConnectionDetailsTable(s.db)
	if err != nil {
		return err
	}

	err = createChannelTable(s.db)
	if err != nil {
		return err
	}

	err = createAttendanceTable(s.db)
	if err != nil {
		return err
	}

	err = createMessageTable(s.db)
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
		name TEXT NOT NULL,
		password TEXT NOT NULL,
		last_login DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)

	return err
}

func createConnectionDetailsTable(db *sql.DB) error {
	// Create the connection_details table if it doesn't exist
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS connection_details (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		host_unique_id TEXT UNIQUE,
		client_unique_id TEXT UNIQUE,
		encryption_key TEXT NOT NULL,
		decryption_key TEXT NOT NULL,
		key_derivation_salt TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)

	if err != nil {
		return err
	}

	// Create an index on the host_unique_id and client_unique_id
	_, err = db.Exec(`
	CREATE UNIQUE INDEX IF NOT EXISTS uniq_connection_details_host_client ON connection_details (host_unique_id, client_unique_id)`)

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

func createAttendanceTable(db *sql.DB) error {
	// Create the attendance table if it doesn't exist
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS attendances (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		channel_id INTEGER,
		joined_at DATETIME,
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (channel_id) REFERENCES channels(id)
	)`)

	if err != nil {
		return err
	}

	// Create an index on the user_id and channel_id columns for faster lookups
	_, err = db.Exec(`
	CREATE INDEX IF NOT EXISTS idx_attendance_user_channel ON attendances (user_id, channel_id)`)

	return err
}

func createMessageTable(db *sql.DB) error {
	// Create the messages table if it doesn't exist
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		channel_id INTEGER,
		content TEXT NOT NULL,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (channel_id) REFERENCES channels(id)
	)`)

	return err
}
