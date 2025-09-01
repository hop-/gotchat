package storage

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/hop-/gotchat/internal/core"
)

const maxRetries = 5

func execWithRetry(s *sql.DB, query string, args ...any) (sql.Result, error) {
	for i := 0; i < maxRetries; i++ {
		result, err := s.Exec(query, args...)
		if err == nil {
			return result, nil
		}

		// Check if the error is a database lock error
		if strings.Contains(err.Error(), "database is locked") || strings.Contains(err.Error(), "busy") {
			// Wait before retrying
			time.Sleep(time.Millisecond * 70 * time.Duration(i+1))
			continue
		}

		return result, err
	}

	return nil, fmt.Errorf("max retries reached")
}

func queryWithRetry(s *sql.DB, query string, args ...any) (*sql.Rows, error) {
	for i := 0; i < maxRetries; i++ {
		rows, err := s.Query(query, args...)
		if err == nil {
			return rows, nil
		}

		// Check if the error is a database lock error
		if strings.Contains(err.Error(), "database is locked") || strings.Contains(err.Error(), "busy") {
			// Wait before retrying
			time.Sleep(time.Millisecond * 70 * time.Duration(i+1))
			continue
		}

		return nil, err
	}

	return nil, fmt.Errorf("max retries reached")
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
