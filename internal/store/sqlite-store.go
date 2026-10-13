package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
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
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		time DATETIME,
		level VARCHAR(255),
		message VARCHAR(255),
		service VARCHAR(255),
		extras JSONB,

		PRIMARY KEY (id AUTOINCREMENT)
	);

	CREATE INDEX IF NOT EXISTS log_time
	ON logs (time);
	`); err != nil {
		return nil, err
	}

	return &SQLiteStore{
		db: db,
	}, nil
}

func prepareReadLogsSQLQuery(filters filter.Filters) (string, []any) {
	wheres := make([]string, 0, 4)
	params := make([]any, 0, 4)

	if !filters.StartDate.IsZero() {
		wheres = append(wheres, "time > ?")
		params = append(params, filters.StartDate)
	}
	if !filters.EndDate.IsZero() {
		wheres = append(wheres, "time < ?")
		params = append(params, filters.EndDate)
	}
	if filters.Level != "" {
		wheres = append(wheres, "level = ?")
		params = append(params, filters.Level)
	}
	if filters.Service != "" {
		wheres = append(wheres, "service = ?")
		params = append(params, filters.Service)
	}

	sqlStr := "SELECT time, level, message, service, extras FROM logs"
	if len(wheres) > 0 {
		sqlStr += " WHERE " + strings.Join(wheres, " AND ")
	}
	return sqlStr, params
}

func prepareInsertLogsSQLQuery(logs []*parser.Log) (string, []any) {
	var sqlStrBuilder strings.Builder
	sqlStrBuilder.WriteString("INSERT INTO logs (time, level, message, service, extras) VALUES ")
	args := make([]any, 0, len(logs)*5)
	for i, log := range logs {
		if i > 0 {
			sqlStrBuilder.WriteString(",")
		}
		sqlStrBuilder.WriteString("(?, ?, ?, ?, ?)")
		args = append(args, log.Time, log.Level, log.Message, log.Service, fmt.Sprintf("%v", log.Extra))
	}

	return sqlStrBuilder.String(), args
}

func (s *SQLiteStore) RetrieveLogs(ctx context.Context, filters filter.Filters) ([]*parser.Log, []error, error) {
	conn, err := s.db.Begin()
	defer conn.Commit()
	if err != nil {
		return nil, nil, err
	}

	query, args := prepareReadLogsSQLQuery(filters)
	stmt, err := conn.PrepareContext(ctx, query)
	if err != nil {
		return nil, nil, err
	}

	rows, err := stmt.QueryContext(ctx, args...)
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

		var extras any
		// TODO: actually return extras
		if err := rows.Scan(&log.Time, &log.Level, &log.Message, &log.Service, &extras); err != nil {
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

// FIXME:  I think this function should get the interface []*parser.Log instead of the channel
// All channel to slice with batch gestion should be in a dedicated function, and this one should only handle session with database life
func (s *SQLiteStore) InsertLogs(ctx context.Context, logChan <-chan *parser.Log, _ <-chan error) ([]*parser.Log, error) {
	logs := make([]*parser.Log, 0, 100)
	shouldInsert := false
	shouldEnd := false
	for {
		select {
		case <-ctx.Done():
			shouldEnd = true
			shouldInsert = len(logs) > 0
		case log, ok := <-logChan:
			if ok {
				logs = append(logs, log)
			} else {
				shouldEnd = true
			}
			shouldInsert = (len(logs) == cap(logs)) || (shouldEnd && len(logs) > 0)
		case <-time.After(500 * time.Millisecond):
			shouldInsert = len(logs) > 0
		}

		if !shouldInsert {
			if shouldEnd {
				return logs, nil
			}
			continue
		}

		conn, err := s.db.Begin()
		if err != nil {
			return nil, err
		}

		query, args := prepareInsertLogsSQLQuery(logs)
		// TODO: for ExecContext and Commit, maybe think about some retry politics or something
		if _, err := conn.ExecContext(ctx, query, args...); err != nil {
			return nil, err
		}

		if err := conn.Commit(); err != nil {
			return nil, err
		}

		logs = []*parser.Log{}
		shouldInsert = false
		if shouldEnd {
			return logs, nil
		}
	}
}
