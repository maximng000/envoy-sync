package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envoy-sync/internal/envfile"
)

var redactCmd = &cobra.Command{
	Use:   "redact <file>",
	Short: "Redact secret values in a .env file and print the result",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]
		modeStr, _ := cmd.Flags().GetString("mode")
		keys, _ := cmd.Flags().GetStringSlice("keys")
		outFile, _ := cmd.Flags().GetString("out")

		mode := envfile.RedactMode(modeStr)
		switch mode {
		case envfile.RedactModeMask, envfile.RedactModeBlank, envfile.RedactModePlaceholder:
			// valid
		default:
			return fmt.Errorf("unknown mode %q: choose mask, blank, or placeholder", modeStr)
		}

		entries, err := envfile.Parse(path)
		if err != nil {
			return fmt.Errorf("failed to parse %s: %w", path, err)
		}

		result := envfile.Redact(entries, mode, keys)

		var sb strings.Builder
		for _, e := range result.Entries {
			sb.WriteString(fmt.Sprintf("%s=%s\n", e.Key, e.Value))
		}

		if outFile != "" {
			if err := os.WriteFile(outFile, []byte(sb.String()), 0644); err != nil {
				return fmt.Errorf("failed to write output: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Redacted %d key(s) → %s\n", len(result.Redacted), outFile)
		} else {
			fmt.Fprint(cmd.OutOrStdout(), sb.String())
		}

		if len(result.Redacted) > 0 {
			fmt.Fprintf(cmd.ErrOrStderr(), "redacted keys: %s\n", strings.Join(result.Redacted, ", "))
		}
		return nil
	},
}

func init() {
	redactCmd.Flags().String("mode", "mask", "Redaction mode: mask | blank | placeholder")
	redactCmd.Flags().StringSlice("keys", nil, "Explicit keys to redact (overrides auto-detection)")
	redactCmd.Flags().String("out", "", "Write output to file instead of stdout")
	rootCmd.AddCommand(redactCmd)
}
