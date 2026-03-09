package storage

import (
	"context"
	"database/sql"

	"github.com/yashnaiduu/Litelog/models"
	_ "modernc.org/sqlite"
)

type Store struct {
	DB *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		DB: db,
	}
}

func InitDB(dbPath string) (*Store, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	pragmas := []string{
		"PRAGMA journal_mode = WAL;",
		"PRAGMA synchronous = NORMAL;",
		"PRAGMA temp_store = MEMORY;",
		"PRAGMA busy_timeout = 5000;",
	}
	for _, pragma := range pragmas {
		if _, err := db.Exec(pragma); err != nil {
			return nil, err
		}
	}

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS logs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
		level TEXT,
		service TEXT,
		message TEXT
	);
	CREATE INDEX IF NOT EXISTS idx_logs_level ON logs(level);
	CREATE INDEX IF NOT EXISTS idx_logs_service ON logs(service);
	CREATE INDEX IF NOT EXISTS idx_logs_timestamp ON logs(timestamp);
	`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		return nil, err
	}

	return NewStore(db), nil
}

func (s *Store) InsertLog(ctx context.Context, level, service, message string) error {
	query := `INSERT INTO logs (level, service, message) VALUES (?, ?, ?)`
	_, err := s.DB.ExecContext(ctx, query, level, service, message)
	return err
}

func (s *Store) InsertLogBatch(ctx context.Context, logs []models.LogEntry) error {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	stmt, err := tx.PrepareContext(ctx, "INSERT INTO logs (level, service, message) VALUES (?, ?, ?)")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, req := range logs {
		if _, err := stmt.ExecContext(ctx, req.Level, req.Service, req.Message); err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (s *Store) QueryLogs(ctx context.Context, level, service string, limit int) ([]models.LogEntry, error) {
	query := "SELECT id, timestamp, level, service, message FROM logs WHERE 1=1"
	args := []interface{}{}

	if level != "" {
		query += " AND level = ?"
		args = append(args, level)
	}
	if service != "" {
		query += " AND service = ?"
		args = append(args, service)
	}

	if limit > 0 {
		query += " ORDER BY timestamp DESC LIMIT ?"
		args = append(args, limit)
	} else {
		query += " ORDER BY timestamp DESC"
	}

	rows, err := s.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []models.LogEntry
	for rows.Next() {
		var log models.LogEntry
		if err := rows.Scan(&log.ID, &log.Timestamp, &log.Level, &log.Service, &log.Message); err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}
	return logs, nil
}

func (s *Store) DeleteOldLogs(ctx context.Context, cutoff string) (int64, error) {
	query := `DELETE FROM logs WHERE timestamp < ?`
	res, err := s.DB.ExecContext(ctx, query, cutoff)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
