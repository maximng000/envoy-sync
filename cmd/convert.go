package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-sync/internal/envfile"
)

var (
	convertFormat string
	convertMask   bool
)

var convertCmd = &cobra.Command{
	Use:   "convert [file]",
	Short: "Convert a .env file to another format",
	Long:  "Convert a .env file to docker-compose, shell export, or dotenv format.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]

		entries, err := envfile.Parse(path)
		if err != nil {
			return fmt.Errorf("failed to parse %s: %w", path, err)
		}

		fmt := envfile.ConvertFormat(convertFormat)
		result, err := envfile.Convert(entries, fmt, convertMask)
		if err != nil {
			return err
		}

		_, err = os.Stdout.WriteString(result.Content)
		return err
	},
}

func init() {
	convertCmd.Flags().StringVarP(
		&convertFormat, "format", "f", "dotenv",
		"Output format: dotenv, shell, docker-compose",
	)
	convertCmd.Flags().BoolVar(
		&convertMask, "mask", false,
		"Mask secret values in output",
	)
	rootCmd.AddCommand(convertCmd)
}
