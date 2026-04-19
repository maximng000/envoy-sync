package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-sync/internal/envfile"
)

var (
	diffMask bool
)

var diffCmd = &cobra.Command{
	Use:   "diff <base> <other>",
	Short: "Diff two .env files and show added, removed, or changed keys",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		basePath := args[0]
		otherPath := args[1]

		base, err := envfile.Parse(basePath)
		if err != nil {
			return fmt.Errorf("reading base file: %w", err)
		}

		other, err := envfile.Parse(otherPath)
		if err != nil {
			return fmt.Errorf("reading other file: %w", err)
		}

		results := envfile.Diff(base, other)

		if len(results) == 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "No differences found.")
			return nil
		}

		for _, d := range results {
			switch d.Status {
			case "added":
				fmt.Fprintf(cmd.OutOrStdout(), "+ %s=%s\n", d.Key, d.OtherValue)
			case "removed":
				fmt.Fprintf(cmd.OutOrStdout(), "- %s=%s\n", d.Key, d.BaseValue)
			case "changed":
				fmt.Fprintf(cmd.OutOrStdout(), "~ %s: %s -> %s\n", d.Key, d.BaseValue, d.OtherValue)
			}
		}

		return nil
	},
}

func init() {
	diffCmd.Flags().BoolVar(&diffMask, "mask", true, "Mask secret values in output")
	if err := rootCmd.RegisterFlagCompletionFunc("mask", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"true", "false"}, cobra.ShellCompDirectiveNoFileComp
	}); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	rootCmd.AddCommand(diffCmd)
}
