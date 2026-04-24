package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envoy-sync/internal/envfile"
)

var cloneCmd = &cobra.Command{
	Use:   "clone <src> <dst>",
	Short: "Clone entries from one .env file into another",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		srcPath, dstPath := args[0], args[1]

		overwrite, _ := cmd.Flags().GetBool("overwrite")
		mask, _ := cmd.Flags().GetBool("mask-secrets")
		filterRaw, _ := cmd.Flags().GetString("keys")

		var filterKeys []string
		if filterRaw != "" {
			for _, k := range strings.Split(filterRaw, ",") {
				if k = strings.TrimSpace(k); k != "" {
					filterKeys = append(filterKeys, k)
				}
			}
		}

		srcEntries, err := envfile.Parse(srcPath)
		if err != nil {
			return fmt.Errorf("reading src: %w", err)
		}

		var dstEntries map[string]string
		if _, statErr := os.Stat(dstPath); os.IsNotExist(statErr) {
			dstEntries = map[string]string{}
		} else {
			dstEntries, err = envfile.Parse(dstPath)
			if err != nil {
				return fmt.Errorf("reading dst: %w", err)
			}
		}

		out, results := envfile.Clone(srcEntries, dstEntries, envfile.CloneOptions{
			FilterKeys:  filterKeys,
			Overwrite:   overwrite,
			MaskSecrets: mask,
		})

		for _, r := range results {
			if r.Skipped {
				fmt.Fprintf(cmd.OutOrStdout(), "SKIP  %s: %s\n", r.Key, r.Reason)
			} else {
				fmt.Fprintf(cmd.OutOrStdout(), "CLONE %s\n", r.Key)
			}
		}

		return envfile.Export(out, dstPath, "dotenv", false)
	},
}

func init() {
	cloneCmd.Flags().Bool("overwrite", false, "Overwrite existing keys in destination")
	cloneCmd.Flags().Bool("mask-secrets", false, "Mask secret values in output log")
	cloneCmd.Flags().String("keys", "", "Comma-separated list of keys to clone (default: all)")
	rootCmd.AddCommand(cloneCmd)
}
