package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	dbPath  string
	rootCmd = &cobra.Command{
		Use:   "pipelinepal",
		Short: "PipelinePal - terminal CRM for real estate",
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&dbPath, "db", defaultDBPath(), "path to sqlite db file")
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(leadCmd)
	rootCmd.AddCommand(followupCmd)
}

func defaultDBPath() string {
	home, _ := os.UserHomeDir()
	return home + "/.local/share/pipelinepal/pipelinepal.db"
}
