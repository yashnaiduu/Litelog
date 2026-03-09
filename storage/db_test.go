package storage

import (
	"testing"
	"time"

	"github.com/yashnaiduu/Litelog/models"
)

func TestDBInitialization(t *testing.T) {
	dbPath := ":memory:" // Use in-memory db

	err := InitDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	if DB == nil {
		t.Fatal("DB instance is nil after successful initialization")
	}

	// Clean up after tests
	defer DB.Close()
}

func TestInsertLog(t *testing.T) {
	dbPath := ":memory:"
	if err := InitDB(dbPath); err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	defer DB.Close()

	err := InsertLog("INFO", "test-service", "this is a test message")
	if err != nil {
		t.Errorf("InsertLog failed: %v", err)
	}

	// Verify insert
	logs, err := QueryLogs("INFO", "test-service", 10)
	if err != nil {
		t.Fatalf("QueryLogs failed: %v", err)
	}

	if len(logs) != 1 {
		t.Errorf("Expected 1 log, got %d", len(logs))
	} else {
		if logs[0].Message != "this is a test message" {
			t.Errorf("Expected message 'this is a test message', got '%s'", logs[0].Message)
		}
	}
}

func TestInsertLogBatch(t *testing.T) {
	dbPath := ":memory:"
	if err := InitDB(dbPath); err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	defer DB.Close()

	batch := []models.LogEntry{
		{Level: "ERROR", Service: "auth", Message: "login failed"},
		{Level: "WARN", Service: "auth", Message: "retry 1"},
		{Level: "ERROR", Service: "payment", Message: "timeout"},
	}

	err := InsertLogBatch(batch)
	if err != nil {
		t.Errorf("InsertLogBatch failed: %v", err)
	}

	logs, err := QueryLogs("ERROR", "", 10) // Should get 2 errors
	if err != nil {
		t.Fatalf("QueryLogs failed: %v", err)
	}

	if len(logs) != 2 {
		t.Errorf("Expected 2 ERROR logs from batch, got %d", len(logs))
	}
}

func TestQueryLogsFiltering(t *testing.T) {
	dbPath := ":memory:"
	if err := InitDB(dbPath); err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	defer DB.Close()

	_ = InsertLog("INFO", "srv-A", "msga1")
	_ = InsertLog("WARN", "srv-B", "msgb1")
	_ = InsertLog("INFO", "srv-C", "msgc1")

	// Query just INFO
	logs, err := QueryLogs("INFO", "", 10)
	if err != nil {
		t.Fatalf("QueryLogs failed: %v", err)
	}
	if len(logs) != 2 {
		t.Errorf("Expected 2 INFO logs, got %d", len(logs))
	}

	// Query specific service
	logs, err = QueryLogs("", "srv-B", 10)
	if err != nil {
		t.Fatalf("QueryLogs failed: %v", err)
	}
	if len(logs) != 1 {
		t.Errorf("Expected 1 log for srv-B, got %d", len(logs))
	}
	if logs[0].Level != "WARN" {
		t.Errorf("Expected WARN level for srv-B log")
	}
}

func TestDeleteOldLogs(t *testing.T) {
	dbPath := ":memory:"
	if err := InitDB(dbPath); err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	defer DB.Close()

	_ = InsertLog("INFO", "test", "old")

	// Dummy future timestamp
	futureCutoff := time.Now().Add(24 * time.Hour).Format("2006-01-02 15:04:05")

	deletedCount, err := DeleteOldLogs(futureCutoff)
	if err != nil {
		t.Fatalf("DeleteOldLogs failed: %v", err)
	}

	if deletedCount != 1 {
		t.Errorf("Expected to delete 1 log, but deleted %d", deletedCount)
	}
}
