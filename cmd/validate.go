package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-sync/internal/envfile"
)

var validateCmd = &cobra.Command{
	Use:   "validate <env-file> [schema-file]",
	Short: "Validate an .env file, optionally against a schema (.env.example)",
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		envPath := args[0]

		env, err := envfile.Parse(envPath)
		if err != nil {
			return fmt.Errorf("parsing env file: %w", err)
		}

		var result envfile.ValidationResult

		if len(args) == 2 {
			schema, err := envfile.Parse(args[1])
			if err != nil {
				return fmt.Errorf("parsing schema file: %w", err)
			}
			result = envfile.ValidateAgainstSchema(env, schema)
		} else {
			result = envfile.ValidateKeys(env)
		}

		if result.OK() {
			fmt.Println("✓ Validation passed")
			return nil
		}

		fmt.Fprintf(os.Stderr, "Validation failed with %d error(s):\n", len(result.Errors))
		for _, e := range result.Errors {
			fmt.Fprintf(os.Stderr, "  - %s\n", e.Error())
		}
		os.Exit(1)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)
}
