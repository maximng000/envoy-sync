package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envoy-sync/internal/envfile"
)

var statsCmd = &cobra.Command{
	Use:   "stats [file]",
	Short: "Show statistics about an env file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]

		entries, err := envfile.Parse(path)
		if err != nil {
			return fmt.Errorf("failed to parse %s: %w", path, err)
		}

		s := envfile.GatherStats(entries)

		fmt.Fprintf(os.Stdout, "File:        %s\n", path)
		fmt.Fprintf(os.Stdout, "Total:       %d\n", s.Total)
		fmt.Fprintf(os.Stdout, "Secrets:     %d\n", s.Secrets)
		fmt.Fprintf(os.Stdout, "Non-secrets: %d\n", s.NonSecrets)
		fmt.Fprintf(os.Stdout, "Empty:       %d\n", s.Empty)
		fmt.Fprintf(os.Stdout, "Unique:      %d\n", s.Unique)
		fmt.Fprintf(os.Stdout, "Duplicates:  %d\n", s.Duplicates)

		topN, _ := cmd.Flags().GetInt("top-prefixes")
		if topN > 0 && len(s.Prefixes) > 0 {
			top := envfile.TopPrefixes(s, topN)
			fmt.Fprintf(os.Stdout, "Top prefixes: %s\n", strings.Join(top, ", "))
		}

		return nil
	},
}

func init() {
	statsCmd.Flags().Int("top-prefixes", 3, "Number of top prefixes to display")
	rootCmd.AddCommand(statsCmd)
}
