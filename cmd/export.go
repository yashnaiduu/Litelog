package cmd

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/yashnaidu/litelog/models"
	"github.com/yashnaidu/litelog/storage"
)

var exportFormat string
var exportService string

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export logs to a file",
	Run: func(cmd *cobra.Command, args []string) {
		if err := storage.InitDB(dbPath); err != nil {
			log.Fatalf("Failed to initialize database: %v", err)
		}

		query := "SELECT id, timestamp, level, service, message FROM logs"
		var queryArgs []interface{}

		if exportService != "" {
			query += " WHERE service = ?"
			queryArgs = append(queryArgs, exportService)
		}
		
		query += " ORDER BY id ASC"

		rows, err := storage.DB.Query(query, queryArgs...)
		if err != nil {
			log.Fatalf("Export query failed: %v", err)
		}
		defer rows.Close()

		var logs []models.LogEntry
		for rows.Next() {
			var entry models.LogEntry
			var ts string
			if err := rows.Scan(&entry.ID, &ts, &entry.Level, &entry.Service, &entry.Message); err != nil {
				log.Printf("Failed to scan row: %v", err)
				continue
			}
			
			logs = append(logs, entry)
		}

		if exportFormat == "csv" {
			writer := csv.NewWriter(os.Stdout)
			writer.Write([]string{"id", "timestamp", "level", "service", "message"})
			for _, entry := range logs {
				writer.Write([]string{
					fmt.Sprintf("%d", entry.ID),
					entry.Timestamp.Format("2006-01-02T15:04:05Z"),
					entry.Level,
					entry.Service,
					entry.Message,
				})
			}
			writer.Flush()
		} else {
			// default to JSON
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", "  ")
			if err := encoder.Encode(logs); err != nil {
				log.Fatalf("Failed to encode JSON: %v", err)
			}
		}
	},
}

func init() {
	exportCmd.Flags().StringVar(&dbPath, "db", "litelog.db", "Path to SQLite database")
	exportCmd.Flags().StringVar(&exportFormat, "format", "json", "Export format (json, csv)")
	exportCmd.Flags().StringVar(&exportService, "service", "", "Filter by service name")
	rootCmd.AddCommand(exportCmd)
}
