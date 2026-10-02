package store

import (
	"database/sql"
	"time"

	"github.com/HugoDrl/zebra/internal/filter"
	"github.com/HugoDrl/zebra/internal/parser"
	_ "github.com/mattn/go-sqlite3"
)

type SQLiteStore struct {
	db *sql.DB
}

func NewSQLiteStore(dbFile string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	if _, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS logs (
		id INTEGER,
		name VARCHAR(255),

		PRIMARY KEY (id AUTOINCREMENT)
	)
	`); err != nil {
		return nil, err
	}

	return &SQLiteStore{
		db: db,
	}, nil
}

func (s *SQLiteStore) RetrieveLogs(filters filter.Filters) ([]*parser.Log, []error, error) {
	conn, err := s.db.Begin()
	if err != nil {
		return nil, nil, err
	}

	stmt, err := conn.Prepare(`SELECT name FROM logs`)
	if err != nil {
		return nil, nil, err
	}

	rows, err := stmt.Query()
	if err != nil {
		return nil, nil, err
	}

	logs := make([]*parser.Log, 0)
	for {
		var log parser.Log
		isEnded := !rows.Next()
		if isEnded {
			break
		}

		if err := rows.Scan(&log.Service); err != nil {
			return nil, nil, err
		}
		logs = append(logs, &log)
	}
	if err := rows.Close(); err != nil {
		return nil, nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, nil, err
	}

	return logs, nil, nil
}

func (s *SQLiteStore) InsertLogs(logChan <-chan *parser.Log, _ <-chan error) ([]*parser.Log, error) {
	logs := make([]*parser.Log, 0, 10)
	shouldInsert := false
	for {
		select {
		case log, ok := <-logChan:
			if !ok {
				return logs, nil
			}
			logs = append(logs, log)
			if len(logs) == cap(logs) {
				shouldInsert = true
			}
		case <-time.After(500 * time.Millisecond):
			if len(logs) > 0 {
				shouldInsert = true
			} else {
			}
		}

		if shouldInsert {
			conn, err := s.db.Begin()
			if err != nil {
				return nil, err
			}

			if _, err := conn.Exec(`INSERT INTO logs (name) VALUES ($1)`, logs[0].Service); err != nil {
				return nil, err
			}
			conn.Commit()
			shouldInsert = false
		}
	}
}
