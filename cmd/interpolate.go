package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-sync/internal/envfile"
)

var (
	interpolateFailOnMissing bool
	interpolateMask         bool
)

var interpolateCmd = &cobra.Command{
	Use:   "interpolate <file>",
	Short: "Resolve variable references within a .env file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]

		entries, err := envfile.Parse(path)
		if err != nil {
			return fmt.Errorf("parsing %q: %w", path, err)
		}

		opts := envfile.InterpolateOptions{
			FailOnMissing: interpolateFailOnMissing,
		}

		resolved, err := envfile.Interpolate(entries, opts)
		if err != nil {
			return fmt.Errorf("interpolation failed: %w", err)
		}

		for _, e := range resolved {
			val := e.Value
			if interpolateMask && envfile.IsSecret(e.Key) {
				val = "****"
			}
			fmt.Fprintf(cmd.OutOrStdout(), "%s=%s\n", e.Key, val)
		}

		return nil
	},
}

func init() {
	interpolateCmd.Flags().BoolVar(&interpolateFailOnMissing, "fail-on-missing", false, "Return an error if a referenced variable is undefined")
	interpolateCmd.Flags().BoolVar(&interpolateMask, "mask", false, "Mask secret values in output")
	rootCmd.AddCommand(interpolateCmd)
	_ = os.Stderr
}
