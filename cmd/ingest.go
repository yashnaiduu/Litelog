package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var ingestUrl string

var ingestCmd = &cobra.Command{
	Use:   "ingest",
	Short: "Pipe logs from stdin to the ingestion server",
	Run: func(cmd *cobra.Command, args []string) {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := scanner.Text()
			
			// Simple heuristic to extract log level
			level := "INFO"
			upperLine := strings.ToUpper(line)
			if strings.Contains(upperLine, "ERROR") || strings.Contains(upperLine, "ERR") {
				level = "ERROR"
			} else if strings.Contains(upperLine, "WARN") {
				level = "WARN"
			} else if strings.Contains(upperLine, "DEBUG") {
				level = "DEBUG"
			}

			reqBody, _ := json.Marshal(map[string]string{
				"level":   level,
				"service": "stdin",
				"message": line,
			})

			resp, err := http.Post(ingestUrl+"/ingest", "application/json", bytes.NewBuffer(reqBody))
			if err != nil {
				log.Printf("Failed to ingest log: %v", err)
				continue
			}
			resp.Body.Close()
		}
		if err := scanner.Err(); err != nil {
			log.Fatalf("Error reading standard input: %v", err)
		}
	},
}

func init() {
	ingestCmd.Flags().StringVar(&ingestUrl, "url", "http://localhost:8080", "Ingestion server URL")
	rootCmd.AddCommand(ingestCmd)
}
