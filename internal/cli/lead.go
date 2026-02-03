package cli

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/yourname/pipelinepal/internal/db"
	"github.com/yourname/pipelinepal/internal/models"
	"github.com/yourname/pipelinepal/internal/repo"
)

var leadCmd = &cobra.Command{
	Use:   "lead",
	Short: "Manage leads/buyers/sellers",
}

var (
	leadName   string
	leadPhone  string
	leadEmail  string
	leadSource string
	leadKind   string
	leadStatus string
	leadNotes  string
	leadFollow string
)

var leadAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a lead (or buyer/seller)",
	RunE: func(cmd *cobra.Command, args []string) error {
		conn, err := db.Open(dbPath)
		if err != nil {
			return err
		}
		defer conn.Close()

		r := repo.NewLeadRepo(conn)

		var next *time.Time
		if leadFollow != "" {
			// Accept: 2026-02-03 or RFC3339
			var t time.Time
			if len(leadFollow) == 10 {
				t, err = time.Parse("2006-01-02", leadFollow)
			} else {
				t, err = time.Parse(time.RFC3339, leadFollow)
			}
			if err != nil {
				return fmt.Errorf("bad --follow date (use YYYY-MM-DD or RFC3339): %w", err)
			}
			next = &t
		}

		l := &models.Lead{
			Name:           leadName,
			Phone:          leadPhone,
			Email:          leadEmail,
			Source:         leadSource,
			Kind:           defaultStr(leadKind, "lead"),
			Status:         defaultStr(leadStatus, "new"),
			Notes:          leadNotes,
			NextFollowUpAt: next,
		}

		if l.Name == "" {
			return fmt.Errorf("--name is required")
		}

		id, err := r.Add(l)
		if err != nil {
			return err
		}

		fmt.Printf("✅ Added %s #%d (%s)\n", l.Kind, id, l.Name)
		return nil
	},
}

var leadListCmd = &cobra.Command{
	Use:   "list",
	Short: "List leads (or buyers/sellers)",
	RunE: func(cmd *cobra.Command, args []string) error {
		conn, err := db.Open(dbPath)
		if err != nil {
			return err
		}
		defer conn.Close()

		r := repo.NewLeadRepo(conn)
		items, err := r.List(leadKind)
		if err != nil {
			return err
		}

		if len(items) == 0 {
			fmt.Println("No records found.")
			return nil
		}

		for _, l := range items {
			fu := ""
			if l.NextFollowUpAt != nil {
				fu = l.NextFollowUpAt.Local().Format("2006-01-02")
			}
			fmt.Printf("#%d %-6s %-10s %-20s follow:%s source:%s\n",
				l.ID, l.Kind, l.Status, l.Name, fu, l.Source)
		}
		return nil
	},
}

func init() {
	leadCmd.AddCommand(leadAddCmd)
	leadCmd.AddCommand(leadListCmd)

	leadAddCmd.Flags().StringVar(&leadName, "name", "", "full name")
	leadAddCmd.Flags().StringVar(&leadPhone, "phone", "", "phone number")
	leadAddCmd.Flags().StringVar(&leadEmail, "email", "", "email")
	leadAddCmd.Flags().StringVar(&leadSource, "source", "", "lead source (referral, open house, online, etc.)")
	leadAddCmd.Flags().StringVar(&leadKind, "kind", "lead", "lead|buyer|seller")
	leadAddCmd.Flags().StringVar(&leadStatus, "status", "new", "new|contacted|nurture|hot|cold|closed|dead")
	leadAddCmd.Flags().StringVar(&leadNotes, "notes", "", "notes")
	leadAddCmd.Flags().StringVar(&leadFollow, "follow", "", "next follow up date (YYYY-MM-DD or RFC3339)")

	leadListCmd.Flags().StringVar(&leadKind, "kind", "", "filter by kind: lead|buyer|seller (empty = all)")
}

func defaultStr(v, d string) string {
	if v == "" {
		return d
	}
	return v
}
