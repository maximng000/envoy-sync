package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-sync/internal/envfile"
)

var (
	scopeName    string
	scopeList    bool
	scopeSummary bool
)

var scopeCmd = &cobra.Command{
	Use:   "scope [file]",
	Short: "Filter or list entries by key scope (prefix)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		entries, err := envfile.Parse(args[0])
		if err != nil {
			return fmt.Errorf("failed to parse file: %w", err)
		}

		if scopeList {
			scopes := envfile.ListScopes(entries)
			if len(scopes) == 0 {
				fmt.Fprintln(os.Stdout, "no scopes detected")
				return nil
			}
			for _, s := range scopes {
				fmt.Fprintln(os.Stdout, s)
			}
			return nil
		}

		if scopeSummary {
			summ := envfile.ScopeSummaryOf(entries)
			fmt.Fprintln(os.Stdout, envfile.FormatScopeSummary(summ))
			return nil
		}

		if scopeName == "" {
			return fmt.Errorf("provide --scope, --list, or --summary")
		}

		result := envfile.Scope(entries, scopeName)
		if len(result.Entries) == 0 {
			fmt.Fprintf(os.Stdout, "no entries found for scope %q\n", result.Scope)
			return nil
		}
		for _, e := range result.Entries {
			val := e.Value
			if envfile.IsSecret(e.Key) {
				val = "***"
			}
			fmt.Fprintf(os.Stdout, "%s=%s\n", e.Key, val)
		}
		return nil
	},
}

func init() {
	scopeCmd.Flags().StringVar(&scopeName, "scope", "", "scope prefix to filter by (e.g. DB, AWS)")
	scopeCmd.Flags().BoolVar(&scopeList, "list", false, "list all detected scopes")
	scopeCmd.Flags().BoolVar(&scopeSummary, "summary", false, "show entry count per scope")
	rootCmd.AddCommand(scopeCmd)
}
