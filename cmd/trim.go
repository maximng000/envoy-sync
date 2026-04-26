package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envoy-sync/internal/envfile"
)

var trimCmd = &cobra.Command{
	Use:   "trim <file>",
	Short: "Remove leading/trailing whitespace from env values",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath := args[0]
		keys, _ := cmd.Flags().GetStringSlice("keys")
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		entries, err := envfile.Parse(filePath)
		if err != nil {
			return fmt.Errorf("failed to parse %s: %w", filePath, err)
		}

		var trimmed map[string]string
		var changes []envfile.TrimResult

		if len(keys) > 0 {
			trimmed, changes = envfile.TrimKeys(entries, keys)
		} else {
			trimmed, changes = envfile.Trim(entries)
		}

		if len(changes) == 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "No values needed trimming.")
			return nil
		}

		for _, c := range changes {
			fmt.Fprintf(cmd.OutOrStdout(), "trimmed: %s = %q -> %q\n", c.Key, c.OldValue, c.NewValue)
		}

		if dryRun {
			fmt.Fprintln(cmd.OutOrStdout(), "(dry-run) no changes written.")
			return nil
		}

		var sb strings.Builder
		for k, v := range trimmed {
			sb.WriteString(fmt.Sprintf("%s=%s\n", k, v))
		}

		if err := os.WriteFile(filePath, []byte(sb.String()), 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", filePath, err)
		}

		fmt.Fprintf(cmd.OutOrStdout(), "wrote %d change(s) to %s\n", len(changes), filePath)
		return nil
	},
}

func init() {
	trimCmd.Flags().StringSlice("keys", nil, "Only trim specific keys (comma-separated)")
	trimCmd.Flags().Bool("dry-run", false, "Preview changes without writing to file")
	rootCmd.AddCommand(trimCmd)
}
