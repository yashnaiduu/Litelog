package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/yashnaidu/litelog/storage"
)


import (
       "encoding/csv"
       "encoding/json"
)

var format string

var queryCmd = &cobra.Command{
       Use:   "query [sql]",
       Short: "Run a SQL query against the logs database",
       Args:  cobra.ExactArgs(1),
       Run: func(cmd *cobra.Command, args []string) {
	       query := args[0]
	       if err := storage.InitDB(dbPath); err != nil {
		       log.Fatalf("Failed to initialize database: %v", err)
	       }

	       rows, err := storage.DB.Query(query)
	       if err != nil {
		       log.Fatalf("Query failed: %v", err)
	       }
	       defer rows.Close()

	       cols, err := rows.Columns()
	       if err != nil {
		       log.Fatalf("Failed to get columns: %v", err)
	       }

	       // Prepare to collect all rows
	       var allRows [][]interface{}
	       var allRowsMap []map[string]interface{}
	       for rows.Next() {
		       values := make([]interface{}, len(cols))
		       valuePtrs := make([]interface{}, len(cols))
		       for i := range cols {
			       valuePtrs[i] = &values[i]
		       }
		       if err := rows.Scan(valuePtrs...); err != nil {
			       log.Fatalf("Failed to scan row: %v", err)
		       }
		       rowCopy := make([]interface{}, len(cols))
		       rowMap := make(map[string]interface{})
		       for i, val := range values {
			       var v interface{}
			       if val == nil {
				       v = nil
			       } else {
				       switch t := val.(type) {
				       case []byte:
					       v = string(t)
				       default:
					       v = t
				       }
			       }
			       rowCopy[i] = v
			       rowMap[cols[i]] = v
		       }
		       allRows = append(allRows, rowCopy)
		       allRowsMap = append(allRowsMap, rowMap)
	       }
	       if err := rows.Err(); err != nil {
		       log.Fatalf("Error iterating over rows: %v", err)
	       }

	       switch format {
	       case "table":
		       w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		       fmt.Fprintln(w, strings.Join(cols, "\t"))
		       separators := make([]string, len(cols))
		       for i, col := range cols {
			       separators[i] = strings.Repeat("-", len(col))
		       }
		       fmt.Fprintln(w, strings.Join(separators, "\t"))
		       for _, row := range allRows {
			       var rowStrs []string
			       for _, val := range row {
				       if val == nil {
					       rowStrs = append(rowStrs, "NULL")
				       } else {
					       rowStrs = append(rowStrs, fmt.Sprintf("%v", val))
				       }
			       }
			       fmt.Fprintln(w, strings.Join(rowStrs, "\t"))
		       }
		       w.Flush()
	       case "json":
		       enc := json.NewEncoder(os.Stdout)
		       enc.SetIndent("", "  ")
		       if err := enc.Encode(allRowsMap); err != nil {
			       log.Fatalf("Failed to encode JSON: %v", err)
		       }
	       case "csv":
		       writer := csv.NewWriter(os.Stdout)
		       if err := writer.Write(cols); err != nil {
			       log.Fatalf("Failed to write CSV header: %v", err)
		       }
		       for _, row := range allRows {
			       var record []string
			       for _, val := range row {
				       if val == nil {
					       record = append(record, "NULL")
				       } else {
					       record = append(record, fmt.Sprintf("%v", val))
				       }
			       }
			       if err := writer.Write(record); err != nil {
				       log.Fatalf("Failed to write CSV row: %v", err)
			       }
		       }
		       writer.Flush()
		       if err := writer.Error(); err != nil {
			       log.Fatalf("CSV writer error: %v", err)
		       }
	       default:
		       log.Fatalf("Unknown format: %s. Supported formats: table, json, csv", format)
	       }
       },
}

func init() {
	queryCmd.Flags().StringVar(&dbPath, "db", "litelog.db", "Path to SQLite database")
	queryCmd.Flags().StringVar(&format, "format", "table", "Output format: table, json, or csv")
	rootCmd.AddCommand(queryCmd)
}
