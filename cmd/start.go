package cmd

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/yashnaiduu/Litelog/server"
	"github.com/yashnaiduu/Litelog/storage"
)

var port string
var dbPath string
var retention string

func parseRetentionToSQLiteModifier(ret string) string {
	if strings.HasSuffix(ret, "d") {
		num := strings.TrimSuffix(ret, "d")
		return "-" + num + " days"
	}
	if strings.HasSuffix(ret, "h") {
		num := strings.TrimSuffix(ret, "h")
		return "-" + num + " hours"
	}
	if strings.HasSuffix(ret, "m") {
		num := strings.TrimSuffix(ret, "m")
		return "-" + num + " minutes"
	}
	return ""
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the LiteLog ingestion server",
	Run: func(cmd *cobra.Command, args []string) {
		if err := storage.InitDB(dbPath); err != nil {
			log.Fatalf("Failed to initialize database: %v", err)
		}

		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		var wg sync.WaitGroup

		if retention != "" {
			modifier := parseRetentionToSQLiteModifier(retention)
			if modifier == "" {
				log.Fatalf("Invalid retention duration format. Use 'd', 'h', or 'm' (e.g. 7d, 24h)")
			}

			numPart := strings.TrimSuffix(strings.TrimSuffix(strings.TrimSuffix(retention, "d"), "h"), "m")
			if _, err := strconv.Atoi(numPart); err != nil {
				log.Fatalf("Invalid retention value: %v", err)
			}

			wg.Add(1)
			go func() {
				defer wg.Done()
				ticker := time.NewTicker(1 * time.Minute)
				defer ticker.Stop()
				for {
					select {
					case <-ctx.Done():
						return
					case <-ticker.C:
						storage.DB.Exec("DELETE FROM logs WHERE timestamp < datetime('now', ?)", modifier)
					}
				}
			}()
			log.Printf("Log retention set to %s\n", retention)
		}

		if err := server.StartHttpServer(ctx, &wg, port); err != nil {
			log.Fatalf("Server failed: %v", err)
		}

		wg.Wait()
		log.Println("LiteLog shutdown complete.")
	},
}

func init() {
	startCmd.Flags().StringVarP(&port, "port", "p", "8080", "Port to listen on")
	startCmd.Flags().StringVar(&dbPath, "db", "litelog.db", "Path to SQLite database")
	startCmd.Flags().StringVar(&retention, "retention", "", "Log retention duration (e.g. 7d, 24h)")
	rootCmd.AddCommand(startCmd)
}
