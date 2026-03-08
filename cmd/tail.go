package cmd

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/yashnaidu/litelog/models"
	"github.com/yashnaidu/litelog/storage"
)

const ansiReset = "\033[0m"

var tailLevel string
var tailService string

var tailCmd = &cobra.Command{
	Use:   "tail",
	Short: "Stream real-time logs",
	Run: func(cmd *cobra.Command, args []string) {
		if err := storage.InitDB(dbPath); err != nil {
			log.Fatalf("Failed to initialize database: %v", err)
		}

		fmt.Println("Tailing logs...")

		var lastID int64 = 0

		row := storage.DB.QueryRow("SELECT COALESCE(MAX(id), 0) FROM logs")
		if err := row.Scan(&lastID); err != nil {
			log.Printf("Failed to get latest log ID: %v", err)
		}

		for {
			query := "SELECT id, timestamp, level, service, message FROM logs WHERE id > ?"
			var queryArgs []interface{}
			queryArgs = append(queryArgs, lastID)

			if tailLevel != "" {
				query += " AND level = ?"
				queryArgs = append(queryArgs, tailLevel)
			}
			if tailService != "" {
				query += " AND service = ?"
				queryArgs = append(queryArgs, tailService)
			}
			query += " ORDER BY id ASC"

			rows, err := storage.DB.Query(query, queryArgs...)
			if err != nil {
				log.Fatalf("Tail query failed: %v", err)
			}

			var maxID int64
			hasNewLogs := false

			for rows.Next() {
				var entry models.LogEntry
				var ts string
				if err := rows.Scan(&entry.ID, &ts, &entry.Level, &entry.Service, &entry.Message); err != nil {
					log.Printf("Failed to scan row: %v", err)
					continue
				}

				if parsedTs, err := time.Parse(time.RFC3339, ts); err == nil {
					ts = parsedTs.Format("15:04:05")
				} else if parsedTs, err := time.Parse("2006-01-02 15:04:05", ts); err == nil {
					ts = parsedTs.Format("15:04:05")
				}

				fmt.Printf("[%s] %-5s %-15s %s\n", ts, colorizeLevel(entry.Level), entry.Service, entry.Message)

				maxID = entry.ID
				hasNewLogs = true
			}
			rows.Close()

			if hasNewLogs {
				lastID = maxID
			}

			time.Sleep(500 * time.Millisecond)
		}
	},
}

func colorizeLevel(level string) string {
	upperLevel := strings.ToUpper(level)

	var color string
	switch upperLevel {
	case "ERROR", "FATAL":
		color = "\033[31m"
	case "WARN", "WARNING":
		color = "\033[33m"
	case "INFO":
		color = "\033[32m"
	case "DEBUG":
		color = "\033[36m"
	default:
		return level
	}

	return color + level + ansiReset
}

func init() {
	tailCmd.Flags().StringVar(&dbPath, "db", "litelog.db", "Path to SQLite database")
	tailCmd.Flags().StringVar(&tailLevel, "level", "", "Filter by log level")
	tailCmd.Flags().StringVar(&tailService, "service", "", "Filter by service name")
	rootCmd.AddCommand(tailCmd)
}
