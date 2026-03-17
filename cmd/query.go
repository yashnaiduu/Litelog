package cmd

import (
	"context"
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
	"github.com/yashnaiduu/Litelog/storage"
)

var queryFormat string

var queryCmd = &cobra.Command{
	Use:   "query [sql]",
	Short: "Run a SQL query against the logs database",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		query := args[0]
		store, err := storage.InitDB(dbPath)
		if err != nil {
			log.Fatalf("Failed to initialize database: %v", err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		rows, err := store.DB.QueryContext(ctx, query)
		if err != nil {
			log.Fatalf("Query failed: %v", err)
		}
		defer rows.Close()

		cols, err := rows.Columns()
		if err != nil {
			log.Fatalf("Failed to get columns: %v", err)
		}

		switch strings.ToLower(queryFormat) {
		case "table":
			if err := writeTable(rows, cols); err != nil {
				log.Fatalf("Failed to write table output: %v", err)
			}
		case "json":
			if err := writeJSON(rows, cols); err != nil {
				log.Fatalf("Failed to write JSON output: %v", err)
			}
		case "csv":
			if err := writeCSV(rows, cols); err != nil {
				log.Fatalf("Failed to write CSV output: %v", err)
			}
		default:
			log.Fatalf("Unknown format %q. Supported formats: table, json, csv", queryFormat)
		}
	},
}

func writeTable(rows *sql.Rows, cols []string) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, strings.Join(cols, "\t"))

	separators := make([]string, len(cols))
	for i, col := range cols {
		separators[i] = strings.Repeat("-", len(col))
	}
	fmt.Fprintln(w, strings.Join(separators, "\t"))

	for rows.Next() {
		values, err := scanRow(rows, len(cols))
		if err != nil {
			return err
		}

		rowStrs := make([]string, 0, len(values))
		for _, val := range values {
			rowStrs = append(rowStrs, formatCell(val))
		}
		fmt.Fprintln(w, strings.Join(rowStrs, "\t"))
	}

	if err := rows.Err(); err != nil {
		return err
	}

	return w.Flush()
}

func writeJSON(rows *sql.Rows, cols []string) error {
	results := make([]map[string]interface{}, 0)

	for rows.Next() {
		values, err := scanRow(rows, len(cols))
		if err != nil {
			return err
		}

		rowMap := make(map[string]interface{}, len(cols))
		for i, col := range cols {
			rowMap[col] = normalizeValue(values[i])
		}
		results = append(results, rowMap)
	}

	if err := rows.Err(); err != nil {
		return err
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(results)
}

func writeCSV(rows *sql.Rows, cols []string) error {
	w := csv.NewWriter(os.Stdout)
	if err := w.Write(cols); err != nil {
		return err
	}

	for rows.Next() {
		values, err := scanRow(rows, len(cols))
		if err != nil {
			return err
		}

		record := make([]string, 0, len(values))
		for _, val := range values {
			record = append(record, formatCell(val))
		}
		if err := w.Write(record); err != nil {
			return err
		}
	}

	if err := rows.Err(); err != nil {
		return err
	}

	w.Flush()
	return w.Error()
}

func scanRow(rows *sql.Rows, colCount int) ([]interface{}, error) {
	values := make([]interface{}, colCount)
	valuePtrs := make([]interface{}, colCount)
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	if err := rows.Scan(valuePtrs...); err != nil {
		return nil, err
	}

	return values, nil
}

func normalizeValue(val interface{}) interface{} {
	switch v := val.(type) {
	case []byte:
		return string(v)
	default:
		return v
	}
}

func formatCell(val interface{}) string {
	if val == nil {
		return "NULL"
	}

	switch v := val.(type) {
	case []byte:
		return string(v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func init() {
	queryCmd.Flags().StringVar(&dbPath, "db", "litelog.db", "Path to SQLite database")
	queryCmd.Flags().StringVar(&queryFormat, "format", "table", "Output format: table, json, or csv")
	rootCmd.AddCommand(queryCmd)
}
