package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envoy-sync/internal/envfile"
)

var (
	resolvePreferEnv  bool
	resolveFailMissing bool
	resolveDefaults   []string
	resolveVerbose    bool
)

var resolveCmd = &cobra.Command{
	Use:   "resolve <file>",
	Short: "Resolve env file values against live environment and defaults",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		entries, err := envfile.Parse(args[0])
		if err != nil {
			return fmt.Errorf("parse: %w", err)
		}

		defaults := parseDefaultsFlag(resolveDefaults)

		opts := envfile.ResolveOptions{
			PreferEnv:     resolvePreferEnv,
			FailOnMissing: resolveFailMissing,
			Defaults:      defaults,
		}

		resolved, err := envfile.Resolve(entries, opts)
		if err != nil {
			return err
		}

		if resolveVerbose {
			fmt.Fprint(os.Stdout, envfile.ResolvedSummary(resolved))
			return nil
		}

		out := envfile.ResolvedToEntries(resolved)
		for _, e := range out {
			fmt.Fprintf(os.Stdout, "%s=%s\n", e.Key, e.Value)
		}
		return nil
	},
}

// parseDefaultsFlag converts ["KEY=val", ...] into a map.
func parseDefaultsFlag(pairs []string) map[string]string {
	if len(pairs) == 0 {
		return nil
	}
	m := make(map[string]string, len(pairs))
	for _, p := range pairs {
		parts := strings.SplitN(p, "=", 2)
		if len(parts) == 2 {
			m[parts[0]] = parts[1]
		}
	}
	return m
}

func init() {
	resolveCmd.Flags().BoolVar(&resolvePreferEnv, "prefer-env", false, "live env vars take precedence over file values")
	resolveCmd.Flags().BoolVar(&resolveFailMissing, "fail-missing", false, "exit with error if any key has no resolved value")
	resolveCmd.Flags().StringArrayVar(&resolveDefaults, "default", nil, "default value for a key: KEY=value (repeatable)")
	resolveCmd.Flags().BoolVar(&resolveVerbose, "verbose", false, "print key, source, and value in tabular form")
	rootCmd.AddCommand(resolveCmd)
}
