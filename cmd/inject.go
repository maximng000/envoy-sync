package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"envoy-sync/internal/envfile"
)

var (
	injectOverwrite bool
	injectKeys      []string
)

var injectCmd = &cobra.Command{
	Use:   "inject <file>",
	Short: "Inject variables from a .env file into the current process environment",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		entries, err := envfile.Parse(args[0])
		if err != nil {
			return fmt.Errorf("failed to parse %s: %w", args[0], err)
		}

		opts := envfile.InjectOptions{
			Overwrite: injectOverwrite,
			Keys:      injectKeys,
		}

		result, err := envfile.Inject(entries, opts)
		if err != nil {
			return err
		}

		if len(result.Injected) > 0 {
			fmt.Fprintf(cmd.OutOrStdout(), "Injected: %s\n", strings.Join(result.Injected, ", "))
		}
		if len(result.Skipped) > 0 {
			fmt.Fprintf(cmd.OutOrStdout(), "Skipped (already set): %s\n", strings.Join(result.Skipped, ", "))
		}
		if len(result.Injected) == 0 && len(result.Skipped) == 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "Nothing to inject.")
		}
		return nil
	},
}

func init() {
	injectCmd.Flags().BoolVar(&injectOverwrite, "overwrite", false, "Overwrite existing environment variables")
	injectCmd.Flags().StringSliceVar(&injectKeys, "keys", nil, "Comma-separated list of keys to inject (default: all)")
	rootCmd.AddCommand(injectCmd)
}
