package store

import (
	"context"
	"database/sql"
	"maps"
	"slices"
	"strings"
	"time"
	"uuid"

	"github.com/HugoDrl/zebra/internal/filter"
	"github.com/HugoDrl/zebra/internal/parser"
	_ "github.com/mattn/go-sqlite3"
)

type SQLiteStore struct {
	db *sql.DB
}

type Log struct {
	Id        uuid.UUID
	CreatedAt time.Time
	Time      time.Time
	Duration  time.Duration
	Level     string
	Message   string
	Service   string
}

type LogExtra struct {
	Id    uuid.UUID
	LogId uuid.UUID
	Name  string
	Value string
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
		id UUID,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		time DATETIME,
		duration DURATION,
		level VARCHAR(255),
		message VARCHAR(255),
		service VARCHAR(255),

		PRIMARY KEY (id)
	);

	CREATE INDEX IF NOT EXISTS log_time
	ON logs (time);

	CREATE TABLE IF NOT EXISTS extras (
		id UUID,
		log_id UUID,
		name VARCHAR(255),
		value VARCHAR(255),

		PRIMARY KEY (id),
		FOREIGN KEY (log_id) REFERENCES log_id,
		UNIQUE(log_id, name)
	)
	`); err != nil {
		return nil, err
	}

	return &SQLiteStore{
		db: db,
	}, nil
}

func prepareLogsToDB(logs []*parser.Log) ([]Log, []LogExtra) {
	dbLogs := make([]Log, len(logs))
	logExtras := make([]LogExtra, 0)
	for i, log := range logs {
		logId := uuid.NewV4()
		dbLog := Log{
			Id:       logId,
			Time:     log.Time,
			Duration: time.Duration(log.Duration),
			Level:    string(log.Level),
			Message:  log.Message,
			Service:  log.Service,
		}
		dbLogs[i] = dbLog

		for name, value := range log.Extra {
			logExtras = append(logExtras, LogExtra{
				Id:    uuid.NewV4(),
				LogId: logId,
				Name:  name,
				Value: value,
			})
		}
	}

	return dbLogs, logExtras
}

func dbLogToLog(log Log) *parser.Log {
	return &parser.Log{
		Time:     log.Time,
		Duration: parser.Duration(log.Duration),
		Level:    parser.Level(log.Level),
		Message:  log.Message,
		Service:  log.Service,
		// extras is initialized to nil and should be updated manually
	}
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

	sqlStr := `
	SELECT
	logs.time,
	logs.level,
	logs.duration,
	logs.message,
	logs.service,
	extras.name,
	extras.value

	FROM logs
	LEFT JOIN extras ON extras.log_id = logs.id`
	if len(wheres) > 0 {
		sqlStr += " WHERE " + strings.Join(wheres, " AND ")
	}
	return sqlStr, params
}

func prepareInsertLogsSQLQuery(logs []Log) (string, []any) {
	var sqlStrBuilder strings.Builder
	sqlStrBuilder.WriteString("INSERT INTO logs (id, time, level, duration, message, service) VALUES ")
	args := make([]any, 0, len(logs)*6)
	for i, log := range logs {
		if i > 0 {
			sqlStrBuilder.WriteString(",")
		}
		sqlStrBuilder.WriteString("(?, ?, ?, ?, ?, ?)")
		args = append(args, log.Id, log.Time, log.Level, log.Duration, log.Message, log.Service)
	}

	return sqlStrBuilder.String(), args
}

func prepareInsertLogsExtrasSQLQuery(extras []LogExtra) (string, []any) {
	var sqlStrBuilder strings.Builder
	sqlStrBuilder.WriteString("INSERT INTO extras (id, log_id, name, value) VALUES ")
	args := make([]any, 0, len(extras)*4)
	for i, extra := range extras {
		if i > 0 {
			sqlStrBuilder.WriteString(",")
		}
		sqlStrBuilder.WriteString("(?, ?, ?, ?)")

		args = append(args, extra.Id, extra.LogId, extra.Name, extra.Value)
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

	logs := make(map[uuid.UUID]*parser.Log)
	for {
		var log Log
		isEnded := !rows.Next()
		if isEnded {
			break
		}

		var extraName, extraValue sql.NullString
		if err := rows.Scan(&log.Time, &log.Level, &log.Duration, &log.Message, &log.Service, &extraName, &extraValue); err != nil {
			return nil, nil, err
		}

		if _, ok := logs[log.Id]; !ok {
			logs[log.Id] = dbLogToLog(log)
		}
		if extraName.Valid {
			if logs[log.Id].Extra == nil {
				logs[log.Id].Extra = make(map[string]string)
			}
			logs[log.Id].Extra[extraName.String] = extraValue.String
		}
	}
	if err := rows.Close(); err != nil {
		return nil, nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, nil, err
	}

	logSlice := slices.Collect(maps.Values(logs))
	if logSlice == nil {
		logSlice = make([]*parser.Log, 0)
	}
	return logSlice, nil, nil
}

// FIXME:  I think this function should get the slice []*parser.Log instead of the channel
// All channel to slice with batch gestion should be in a dedicated function, and this one should only handle session with database life
func (s *SQLiteStore) InsertLogs(ctx context.Context, logChan <-chan *parser.Log, _ <-chan error) ([]*parser.Log, error) {
	logs := make([]*parser.Log, 0, 1000)
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

		dbLogs, extras := prepareLogsToDB(logs)
		logsQuery, logArgs := prepareInsertLogsSQLQuery(dbLogs)
		// TODO: for ExecContext and Commit, maybe think about some retry policies or something
		if _, err = conn.ExecContext(ctx, logsQuery, logArgs...); err != nil {
			return nil, err
		}
		if len(extras) > 0 {
			extrasQuery, extraArgs := prepareInsertLogsExtrasSQLQuery(extras)
			if _, err = conn.ExecContext(ctx, extrasQuery, extraArgs...); err != nil {
				return nil, err
			}
		}

		if err := conn.Commit(); err != nil {
			return nil, err
		}

		logs = logs[:0]
		shouldInsert = false
		if shouldEnd {
			return logs, nil
		}
	}
}
