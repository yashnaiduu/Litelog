package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/yashnaidu/litelog/models"
	"github.com/yashnaidu/litelog/storage"
)

var LogQueue chan models.LogEntry

func init() {
	// Buffer allowing massive bursts, up to 100k
	LogQueue = make(chan models.LogEntry, 100000)
}

func StartAsyncWorker() {
	go func() {
		var batch []models.LogEntry
		ticker := time.NewTicker(200 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case entry := <-LogQueue:
				batch = append(batch, entry)
				if len(batch) >= 100 {
					if err := storage.InsertLogBatch(batch); err != nil {
						log.Printf("Batch insert failed: %v", err)
					}
					batch = batch[:0]
				}
			case <-ticker.C:
				if len(batch) > 0 {
					if err := storage.InsertLogBatch(batch); err != nil {
						log.Printf("Batch insert failed: %v", err)
					}
					batch = batch[:0]
				}
			}
		}
	}()
}

type IngestRequest struct {
	Level   string `json:"level"`
	Service string `json:"service"`
	Message string `json:"message"`
}

func StartHttpServer(port string) error {
	StartAsyncWorker()

	http.HandleFunc("/ingest", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req IngestRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		// Send to async queue immediately
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

	fmt.Printf("LiteLog server listening on port %s\n", port)
	return http.ListenAndServe(":"+port, nil)
}

