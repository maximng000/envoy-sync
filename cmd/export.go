package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-sync/internal/envfile"
)

var (
	exportFormat     string
	exportMaskSecret bool
)

var exportCmd = &cobra.Command{
	Use:   "export [file]",
	Short: "Export an .env file in a specified format",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]

		f, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("opening file: %w", err)
		}
		defer f.Close()

		entries, err := envfile.Parse(f)
		if err != nil {
			return fmt.Errorf("parsing file: %w", err)
		}

		m := make(map[string]string, len(entries))
		for _, e := range entries {
			m[e.Key] = e.Value
		}

		return envfile.Export(os.Stdout, m, envfile.Format(exportFormat), exportMaskSecret)
	},
}

func init() {
	exportCmd.Flags().StringVarP(&exportFormat, "format", "f", "dotenv", "Output format: dotenv, export, json")
	exportCmd.Flags().BoolVarP(&exportMaskSecret, "mask", "m", false, "Mask secret values in output")
	RootCmd.AddCommand(exportCmd)
}
