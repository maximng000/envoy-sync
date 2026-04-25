package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-sync/internal/envfile"
)

var schemaCmd = &cobra.Command{
	Use:   "schema",
	Short: "Validate an .env file against a JSON schema",
	RunE: func(cmd *cobra.Command, args []string) error {
		envPath, _ := cmd.Flags().GetString("file")
		schemaPath, _ := cmd.Flags().GetString("schema")

		data, err := os.ReadFile(envPath)
		if err != nil {
			return fmt.Errorf("read env file: %w", err)
		}
		entries, err := envfile.Parse(string(data))
		if err != nil {
			return fmt.Errorf("parse env file: %w", err)
		}

		schema, err := envfile.LoadSchema(schemaPath)
		if err != nil {
			return fmt.Errorf("load schema: %w", err)
		}

		violations := envfile.CheckSchema(entries, schema)
		if len(violations) == 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "✓ env file satisfies schema")
			return nil
		}

		fmt.Fprintf(cmd.OutOrStdout(), "schema violations (%d):\n", len(violations))
		for _, v := range violations {
			fmt.Fprintf(cmd.OutOrStdout(), "  - %s\n", v)
		}
		return fmt.Errorf("schema validation failed with %d violation(s)", len(violations))
	},
}

func init() {
	schemaCmd.Flags().StringP("file", "f", ".env", "path to the .env file")
	schemaCmd.Flags().StringP("schema", "s", ".env.schema.json", "path to the JSON schema file")
	rootCmd.AddCommand(schemaCmd)
}
