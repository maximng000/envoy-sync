package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-sync/internal/envfile"
)

var lintCmd = &cobra.Command{
	Use:   "lint [file]",
	Short: "Lint an .env file for style and correctness issues",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]

		f, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("could not open file: %w", err)
		}
		defer f.Close()

		entries, err := envfile.Parse(f)
		if err != nil {
			return fmt.Errorf("parse error: %w", err)
		}

		result := envfile.Lint(entries)

		if len(result.Issues) == 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "No issues found.")
			return nil
		}

		for _, issue := range result.Issues {
			fmt.Fprintln(cmd.OutOrStdout(), issue.String())
		}

		if result.HasErrors() {
			return fmt.Errorf("lint failed with errors")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(lintCmd)
}
