package cmd

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/yashnaidu/litelog/server"
	"github.com/yashnaidu/litelog/storage"
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

		if retention != "" {
			modifier := parseRetentionToSQLiteModifier(retention)
			if modifier == "" {
				log.Fatalf("Invalid retention duration format. Use 'd', 'h', or 'm' (e.g. 7d, 24h)")
			}

			// Validate it's a number
			numPart := strings.TrimSuffix(strings.TrimSuffix(strings.TrimSuffix(retention, "d"), "h"), "m")
			if _, err := strconv.Atoi(numPart); err != nil {
				log.Fatalf("Invalid retention value: %v", err)
			}

			go func() {
				for {
					storage.DB.Exec("DELETE FROM logs WHERE timestamp < datetime('now', ?)", modifier)
					time.Sleep(1 * time.Minute)
				}
			}()
			log.Printf("Log retention set to %s\n", retention)
		}

		if err := server.StartHttpServer(port); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	},
}

func init() {
	startCmd.Flags().StringVarP(&port, "port", "p", "8080", "Port to listen on")
	startCmd.Flags().StringVar(&dbPath, "db", "litelog.db", "Path to SQLite database")
	startCmd.Flags().StringVar(&retention, "retention", "", "Log retention duration (e.g. 7d, 24h)")
	rootCmd.AddCommand(startCmd)
}
