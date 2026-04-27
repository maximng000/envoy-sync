package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envoy-sync/internal/envfile"
)

var tagCmd = &cobra.Command{
	Use:   "tag [file]",
	Short: "Tag env entries with a label",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath := args[0]
		tag, _ := cmd.Flags().GetString("tag")
		keysRaw, _ := cmd.Flags().GetString("keys")
		summary, _ := cmd.Flags().GetBool("summary")

		if tag == "" {
			return fmt.Errorf("--tag is required")
		}

		entries, err := envfile.Parse(filePath)
		if err != nil {
			return fmt.Errorf("failed to parse file: %w", err)
		}

		var keys []string
		if keysRaw != "" {
			for _, k := range strings.Split(keysRaw, ",") {
				k = strings.TrimSpace(k)
				if k != "" {
					keys = append(keys, k)
				}
			}
		}

		result, err := envfile.Tag(entries, keys, tag)
		if err != nil {
			return fmt.Errorf("tag failed: %w", err)
		}

		if summary {
			fmt.Print(envfile.TagSummary(result.Tagged))
			return nil
		}

		for _, te := range result.Tagged {
			fmt.Fprintf(os.Stdout, "%s=%s  [%s]\n", te.Key, te.Value, te.Tag)
		}
		if len(result.Skipped) > 0 {
			fmt.Fprintf(os.Stdout, "skipped: %s\n", strings.Join(result.Skipped, ", "))
		}
		return nil
	},
}

func init() {
	tagCmd.Flags().String("tag", "", "Tag label to apply (required)")
	tagCmd.Flags().String("keys", "", "Comma-separated keys to tag (default: all)")
	tagCmd.Flags().Bool("summary", false, "Print grouped summary output")
	rootCmd.AddCommand(tagCmd)
}
