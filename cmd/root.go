package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "litelog",
	Short: "LiteLog — centralized logging without the infrastructure.",
	Long:  `LiteLog — the SQLite of logging systems. A single binary for log ingestion, querying, and real-time streaming.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
