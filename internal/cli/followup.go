package cli

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/yourname/pipelinepal/internal/db"
)

var followupCmd = &cobra.Command{
	Use:   "followups",
	Short: "Show follow-ups due",
}

var followupTodayCmd = &cobra.Command{
	Use:   "today",
	Short: "Show follow-ups due today or overdue",
	RunE: func(cmd *cobra.Command, args []string) error {
		conn, err := db.Open(dbPath)
		if err != nil {
			return err
		}
		defer conn.Close()

		// Compare by date; store RFC3339 strings in sqlite, use date() helper.
		today := time.Now().Format("2006-01-02")

		rows, err := conn.Query(`
			SELECT id, kind, status, name, source, next_follow_up_at
			FROM leads
			WHERE next_follow_up_at IS NOT NULL
			  AND date(next_follow_up_at) <= date(?)
			ORDER BY date(next_follow_up_at) ASC
		`, today)
		if err != nil {
			return err
		}
		defer rows.Close()

		found := false
		for rows.Next() {
			found = true
			var id int64
			var kind, status, name, source, next string
			if err := rows.Scan(&id, &kind, &status, &name, &source, &next); err != nil {
				return err
			}
			fmt.Printf("#%d %-6s %-10s %-20s due:%s source:%s\n", id, kind, status, name, next[:10], source)
		}
		if !found {
			fmt.Println("✅ No follow-ups due today.")
		}
		return rows.Err()
	},
}

func init() {
	followupCmd.AddCommand(followupTodayCmd)
}
