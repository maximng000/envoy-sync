package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"envoy-sync/internal/envfile"
)

var (
	sortStrategy string
	sortFormat   string
)

var sortCmd = &cobra.Command{
	Use:   "sort [file]",
	Short: "Sort entries in a .env file by a given strategy",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		entries, err := envfile.Parse(args[0])
		if err != nil {
			return fmt.Errorf("parse: %w", err)
		}

		res, err := envfile.Sort(entries, envfile.SortStrategy(sortStrategy))
		if err != nil {
			return fmt.Errorf("sort: %w", err)
		}

		out, err := envfile.Export(res.Entries, envfile.ExportOptions{
			Format: sortFormat,
			Mask:   false,
		})
		if err != nil {
			return fmt.Errorf("export: %w", err)
		}

		fmt.Fprint(cmd.OutOrStdout(), out)
		return nil
	},
}

func init() {
	sortCmd.Flags().StringVarP(&sortStrategy, "strategy", "s", "alpha",
		"sort strategy: alpha, alpha-desc, secret, length")
	sortCmd.Flags().StringVarP(&sortFormat, "format", "f", "dotenv",
		"output format: dotenv, export, json")
	rootCmd.AddCommand(sortCmd)
}
