package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"

	"envoy-sync/internal/envfile"
)

var (
	flattenPrefix    string
	flattenDelimiter string
	flattenSummary   bool
)

var flattenCmd = &cobra.Command{
	Use:   "flatten <file>",
	Short: "Group and display env keys by their prefix",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		entries, err := envfile.Parse(args[0])
		if err != nil {
			return fmt.Errorf("failed to parse file: %w", err)
		}

		flat := make(map[string]string, len(entries))
		for _, e := range entries {
			flat[e.Key] = e.Value
		}

		if flattenSummary {
			summary := envfile.FlattenSummary(flat, flattenDelimiter)
			prefixes := make([]string, 0, len(summary))
			for p := range summary {
				prefixes = append(prefixes, p)
			}
			sort.Strings(prefixes)
			for _, p := range prefixes {
				fmt.Fprintf(os.Stdout, "[%s] (%d keys)\n", p, len(summary[p]))
				for _, k := range summary[p] {
					fmt.Fprintf(os.Stdout, "  %s\n", k)
				}
			}
			return nil
		}

		results := envfile.Flatten(flat, flattenDelimiter, flattenPrefix)
		if len(results) == 0 {
			fmt.Fprintln(os.Stdout, "no matching keys found")
			return nil
		}

		currentPrefix := ""
		for _, r := range results {
			if r.Prefix != currentPrefix {
				currentPrefix = r.Prefix
				label := currentPrefix
				if label == "" {
					label = "(no prefix)"
				}
				fmt.Fprintf(os.Stdout, "[%s]\n", label)
			}
			fmt.Fprintf(os.Stdout, "  %s=%s\n", r.Key, r.Value)
		}
		return nil
	},
}

func init() {
	flattenCmd.Flags().StringVarP(&flattenPrefix, "prefix", "p", "", "filter by prefix")
	flattenCmd.Flags().StringVarP(&flattenDelimiter, "delimiter", "d", "_", "key delimiter")
	flattenCmd.Flags().BoolVarP(&flattenSummary, "summary", "s", false, "show summary grouped by prefix")
	rootCmd.AddCommand(flattenCmd)
}
