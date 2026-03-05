package storage

import (
	"database/sql"

	_ "modernc.org/sqlite"
	"github.com/yashnaidu/litelog/models"
)

var DB *sql.DB

func InitDB(dbPath string) error {
	var err error
	DB, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}

	// Enable WAL mode and optimize pragmas for high writes
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

