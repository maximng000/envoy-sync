package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-sync/internal/envfile"
)

var syncOverride bool

var syncCmd = &cobra.Command{
	Use:   "sync <base> <source>",
	Short: "Sync keys from source .env into base .env",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		base, err := envfile.Parse(args[0])
		if err != nil {
			return fmt.Errorf("reading base: %w", err)
		}
		src, err := envfile.Parse(args[1])
		if err != nil {
			return fmt.Errorf("reading source: %w", err)
		}

		strategy := envfile.StrategySkip
		if syncOverride {
			strategy = envfile.StrategyOverride
		}

		_, sr := envfile.Sync(base, src, strategy)

		for _, k := range sr.Applied {
			fmt.Fprintf(os.Stdout, "  applied : %s\n", k)
		}
		for _, k := range sr.Skipped {
			fmt.Fprintf(os.Stdout, "  skipped : %s\n", k)
		}
		for _, c := range sr.Conflicts {
			fmt.Fprintf(os.Stdout, "  conflict: %s\n", c)
		}
		return nil
	},
}

func init() {
	syncCmd.Flags().BoolVarP(&syncOverride, "override", "o", false, "override conflicting keys with source values")
	rootCmd.AddCommand(syncCmd)
}
