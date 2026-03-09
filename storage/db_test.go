package storage

import (
	"context"
	"testing"
	"time"

	"github.com/yashnaiduu/Litelog/models"
)

func TestDBInitialization(t *testing.T) {
	dbPath := ":memory:" // Use in-memory db

	store, err := InitDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	if store.DB == nil {
		t.Fatal("DB instance is nil after successful initialization")
	}

	// Clean up after tests
	defer store.DB.Close()
}

func TestInsertLog(t *testing.T) {
	dbPath := ":memory:"
	store, err := InitDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	defer store.DB.Close()

	err = store.InsertLog(context.Background(), "INFO", "test-service", "this is a test message")
	if err != nil {
		t.Errorf("InsertLog failed: %v", err)
	}

	// Verify insert
	logs, err := store.QueryLogs(context.Background(), "INFO", "test-service", 10)
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
	store, err := InitDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	defer store.DB.Close()

	batch := []models.LogEntry{
		{Level: "ERROR", Service: "auth", Message: "login failed"},
		{Level: "WARN", Service: "auth", Message: "retry 1"},
		{Level: "ERROR", Service: "payment", Message: "timeout"},
	}

	err = store.InsertLogBatch(context.Background(), batch)
	if err != nil {
		t.Errorf("InsertLogBatch failed: %v", err)
	}

	logs, err := store.QueryLogs(context.Background(), "ERROR", "", 10) // Should get 2 errors
	if err != nil {
		t.Fatalf("QueryLogs failed: %v", err)
	}

	if len(logs) != 2 {
		t.Errorf("Expected 2 ERROR logs from batch, got %d", len(logs))
	}
}

func TestQueryLogsFiltering(t *testing.T) {
	dbPath := ":memory:"
	store, err := InitDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	defer store.DB.Close()

	_ = store.InsertLog(context.Background(), "INFO", "srv-A", "msga1")
	_ = store.InsertLog(context.Background(), "WARN", "srv-B", "msgb1")
	_ = store.InsertLog(context.Background(), "INFO", "srv-C", "msgc1")

	// Query just INFO
	logs, err := store.QueryLogs(context.Background(), "INFO", "", 10)
	if err != nil {
		t.Fatalf("QueryLogs failed: %v", err)
	}
	if len(logs) != 2 {
		t.Errorf("Expected 2 INFO logs, got %d", len(logs))
	}

	// Query specific service
	logs, err = store.QueryLogs(context.Background(), "", "srv-B", 10)
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
	store, err := InitDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	defer store.DB.Close()

	_ = store.InsertLog(context.Background(), "INFO", "test", "old")

	// Dummy future timestamp
	futureCutoff := time.Now().Add(24 * time.Hour).Format("2006-01-02 15:04:05")

	deletedCount, err := store.DeleteOldLogs(context.Background(), futureCutoff)
	if err != nil {
		t.Fatalf("DeleteOldLogs failed: %v", err)
	}

	if deletedCount != 1 {
		t.Errorf("Expected to delete 1 log, but deleted %d", deletedCount)
	}
}
