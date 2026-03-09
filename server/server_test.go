package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yashnaiduu/Litelog/models"
	"github.com/yashnaiduu/Litelog/storage"
)

func TestIngestEndpoint(t *testing.T) {
	// Initialize an in-memory database to avoid relying on filesystem db
	if err := storage.InitDB(":memory:"); err != nil {
		t.Fatalf("Failed to init storage DB: %v", err)
	}
	defer storage.DB.Close()

	// Ensure the queue is initialized clean
	LogQueue = make(chan models.LogEntry, 10)

	// Ingest endpoint logic
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req IngestRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		select {
		case LogQueue <- models.LogEntry{
			Level:   req.Level,
			Service: req.Service,
			Message: req.Message,
		}:
		default:
			http.Error(w, "Server overloaded", http.StatusServiceUnavailable)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok\n"))
	})

	// Test case 1: Valid payload
	payload := `{"level":"INFO","service":"test-service","message":"test message"}`
	req, err := http.NewRequest("POST", "/ingest", bytes.NewBuffer([]byte(payload)))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Verify log made it to our channel
	select {
	case entry := <-LogQueue:
		if entry.Message != "test message" {
			t.Errorf("Expected message 'test message', got '%s'", entry.Message)
		}
	case <-time.After(1 * time.Second):
		t.Error("Log queue timeout: failed to receive log entry")
	}

	// Test case 2: Invalid JSON
	reqBad, _ := http.NewRequest("POST", "/ingest", bytes.NewBuffer([]byte(`{invalid-json`)))
	rrBad := httptest.NewRecorder()
	handler.ServeHTTP(rrBad, reqBad)

	if status := rrBad.Code; status != http.StatusBadRequest {
		t.Errorf("Expected Bad Request (400) for invalid JSON, got %v", status)
	}

	// Test case 3: Wrong Method
	reqMethod, _ := http.NewRequest("GET", "/ingest", nil)
	rrMethod := httptest.NewRecorder()
	handler.ServeHTTP(rrMethod, reqMethod)

	if status := rrMethod.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("Expected Method Not Allowed (405) for GET request, got %v", status)
	}
}
