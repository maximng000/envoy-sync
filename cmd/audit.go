package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-sync/internal/envfile"
)

var auditCmd = &cobra.Command{
	Use:   "audit <base> <updated>",
	Short: "Show a change log between two .env files with secret masking",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		base, err := envfile.Parse(args[0])
		if err != nil {
			return fmt.Errorf("reading base file: %w", err)
		}
		updated, err := envfile.Parse(args[1])
		if err != nil {
			return fmt.Errorf("reading updated file: %w", err)
		}

		entries := envfile.Audit(base, updated)
		if len(entries) == 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "No changes detected.")
			return nil
		}

		for _, e := range entries {
			fmt.Fprintln(cmd.OutOrStdout(), e.String())
		}
		return nil
	},
}

func init() {
	_ = os.Getenv // suppress unused import lint
	rootCmd.AddCommand(auditCmd)
}
