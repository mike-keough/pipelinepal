package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yourname/pipelinepal/internal/db"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize database and run migrations",
	RunE: func(cmd *cobra.Command, args []string) error {
		conn, err := db.Open(dbPath)
		if err != nil {
			return err
		}
		defer conn.Close()

		if err := db.Migrate(conn); err != nil {
			return err
		}
		fmt.Println("✅ Database initialized:", dbPath)
		return nil
	},
}
