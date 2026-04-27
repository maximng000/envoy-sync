package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-sync/internal/envfile"
)

var (
	normalizeUppercase  bool
	normalizeTrim       bool
	normalizeRemoveEmpty bool
	normalizeSort       bool
	normalizeWrite      bool
)

var normalizeCmd = &cobra.Command{
	Use:   "normalize [file]",
	Short: "Normalize a .env file (uppercase keys, trim values, remove empty)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]

		entries, err := envfile.Parse(path)
		if err != nil {
			return fmt.Errorf("parse error: %w", err)
		}

		opts := envfile.NormalizeOptions{
			UppercaseKeys: normalizeUppercase,
			TrimValues:    normalizeTrim,
			RemoveEmpty:   normalizeRemoveEmpty,
			SortAlpha:     normalizeSort,
		}

		result := envfile.Normalize(entries, opts)

		if normalizeWrite {
			out, err := envfile.Export(result.Entries, envfile.ExportOptions{Format: "dotenv"})
			if err != nil {
				return fmt.Errorf("export error: %w", err)
			}
			if err := os.WriteFile(path, []byte(out), 0644); err != nil {
				return fmt.Errorf("write error: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Wrote normalized file: %s\n", path)
		}

		fmt.Fprint(cmd.OutOrStdout(), envfile.NormalizeSummary(result))

		for _, k := range result.Modified {
			fmt.Fprintf(cmd.OutOrStdout(), "  modified: %s\n", k)
		}
		for _, k := range result.Removed {
			fmt.Fprintf(cmd.OutOrStdout(), "  removed:  %s\n", k)
		}
		return nil
	},
}

func init() {
	normalizeCmd.Flags().BoolVar(&normalizeUppercase, "uppercase", false, "Convert all keys to uppercase")
	normalizeCmd.Flags().BoolVar(&normalizeTrim, "trim", false, "Trim whitespace from values")
	normalizeCmd.Flags().BoolVar(&normalizeRemoveEmpty, "remove-empty", false, "Remove entries with empty values")
	normalizeCmd.Flags().BoolVar(&normalizeSort, "sort", false, "Sort keys alphabetically")
	normalizeCmd.Flags().BoolVar(&normalizeWrite, "write", false, "Write normalized output back to file")
	rootCmd.AddCommand(normalizeCmd)
}
