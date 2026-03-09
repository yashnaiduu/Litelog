package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/yashnaiduu/Litelog/models"
	"github.com/yashnaiduu/Litelog/storage"
)

var LogQueue chan models.LogEntry

func init() {
	LogQueue = make(chan models.LogEntry, 100000)
}

func StartAsyncWorker(ctx context.Context, wg *sync.WaitGroup, store *storage.Store) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		var batch []models.LogEntry
		ticker := time.NewTicker(200 * time.Millisecond)
		defer ticker.Stop()

		flush := func() {
			if len(batch) > 0 {
				// We intentionally decouple this context from the parent context.
				// If the server receives SIGTERM, the parent context cancels, triggering flush.
				// We don't want the insert to instantly cancel, so we use a fresh 5s timeout.
				insertCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				if err := store.InsertLogBatch(insertCtx, batch); err != nil {
					log.Printf("Batch insert failed: %v", err)
				}
				batch = batch[:0]
			}
		}

		for {
			select {
			case <-ctx.Done():
				// Drain the queue to ensure no logs are lost
				for {
					select {
					case entry := <-LogQueue:
						batch = append(batch, entry)
					default:
						flush()
						return
					}
				}
			case entry := <-LogQueue:
				batch = append(batch, entry)
				if len(batch) >= 100 {
					flush()
				}
			case <-ticker.C:
				flush()
			}
		}
	}()
}

type IngestRequest struct {
	Level   string `json:"level"`
	Service string `json:"service"`
	Message string `json:"message"`
}

func StartHttpServer(ctx context.Context, wg *sync.WaitGroup, port string, store *storage.Store) error {
	StartAsyncWorker(ctx, wg, store)

	mux := http.NewServeMux()
	mux.HandleFunc("/ingest", func(w http.ResponseWriter, r *http.Request) {
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
		_, _ = w.Write([]byte("ok\n"))
	})

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Printf("HTTP server shutdown error: %v", err)
		}
	}()

	fmt.Printf("LiteLog server listening on port %s\n", port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}
