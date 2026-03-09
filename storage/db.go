package storage

import (
	"database/sql"

	"github.com/yashnaiduu/Litelog/models"
	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB(dbPath string) error {
	var err error
	DB, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}

	pragmas := []string{
		"PRAGMA journal_mode = WAL;",
		"PRAGMA synchronous = NORMAL;",
		"PRAGMA temp_store = MEMORY;",
		"PRAGMA busy_timeout = 5000;",
	}
	for _, pragma := range pragmas {
		if _, err := DB.Exec(pragma); err != nil {
			return err
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
	_, err = DB.Exec(createTableQuery)
	return err
}

func InsertLog(level, service, message string) error {
	query := `INSERT INTO logs (level, service, message) VALUES (?, ?, ?)`
	_, err := DB.Exec(query, level, service, message)
	return err
}

func InsertLogBatch(logs []models.LogEntry) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("INSERT INTO logs (level, service, message) VALUES (?, ?, ?)")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, req := range logs {
		if _, err := stmt.Exec(req.Level, req.Service, req.Message); err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func QueryLogs(level, service string, limit int) ([]models.LogEntry, error) {
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

	rows, err := DB.Query(query, args...)
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

func DeleteOldLogs(cutoff string) (int64, error) {
	query := `DELETE FROM logs WHERE timestamp < ?`
	res, err := DB.Exec(query, cutoff)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
