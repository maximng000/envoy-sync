package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envoy-sync/internal/envfile"
)

func init() {
	var style string
	var sortKeys bool
	var maskSecret bool
	var inPlace bool

	cmd := &cobra.Command{
		Use:   "format <file>",
		Short: "Format a .env file with a consistent style",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]

			entries, err := envfile.Parse(path)
			if err != nil {
				return fmt.Errorf("parse error: %w", err)
			}

			var fs envfile.FormatStyle
			switch strings.ToLower(style) {
			case "aligned":
				fs = envfile.StyleAligned
			case "spaced":
				fs = envfile.StyleSpaced
			case "compact", "":
				fs = envfile.StyleCompact
			default:
				return fmt.Errorf("unknown style %q: use compact, spaced, or aligned", style)
			}

			result := envfile.Format(entries, envfile.FormatOptions{
				Style:      fs,
				SortKeys:   sortKeys,
				MaskSecret: maskSecret,
			})

			output := strings.Join(result.Lines, "\n") + "\n"

			if inPlace {
				if err := os.WriteFile(path, []byte(output), 0644); err != nil {
					return fmt.Errorf("write error: %w", err)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "formatted %s (%d lines modified)\n", path, result.Modified)
			} else {
				fmt.Fprint(cmd.OutOrStdout(), output)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&style, "style", "s", "compact", "Format style: compact, spaced, aligned")
	cmd.Flags().BoolVar(&sortKeys, "sort", false, "Sort keys alphabetically")
	cmd.Flags().BoolVar(&maskSecret, "mask", false, "Mask secret values in output")
	cmd.Flags().BoolVarP(&inPlace, "in-place", "i", false, "Write formatted output back to file")

	rootCmd.AddCommand(cmd)
}
