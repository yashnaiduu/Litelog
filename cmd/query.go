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

		// Setup tabwriter for aligned output
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, strings.Join(cols, "\t"))

		// Print separator
		separators := make([]string, len(cols))
		for i, col := range cols {
			separators[i] = strings.Repeat("-", len(col))
		}
		fmt.Fprintln(w, strings.Join(separators, "\t"))

		values := make([]interface{}, len(cols))
		valuePtrs := make([]interface{}, len(cols))
		for i := range cols {
			valuePtrs[i] = &values[i]
		}

		for rows.Next() {
			if err := rows.Scan(valuePtrs...); err != nil {
				log.Fatalf("Failed to scan row: %v", err)
			}

			var rowStrs []string
			for _, val := range values {
				if val == nil {
					rowStrs = append(rowStrs, "NULL")
				} else {
					switch v := val.(type) {
					case []byte:
						rowStrs = append(rowStrs, string(v))
					default:
						rowStrs = append(rowStrs, fmt.Sprintf("%v", v))
					}
				}
			}
			fmt.Fprintln(w, strings.Join(rowStrs, "\t"))
		}

		if err := rows.Err(); err != nil {
			log.Fatalf("Error iterating over rows: %v", err)
		}
		w.Flush()
	},
}

func init() {
	queryCmd.Flags().StringVar(&dbPath, "db", "litelog.db", "Path to SQLite database")
	rootCmd.AddCommand(queryCmd)
}
// placeholder
