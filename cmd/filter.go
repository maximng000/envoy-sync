package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envoy-sync/internal/envfile"
)

var filterCmd = &cobra.Command{
	Use:   "filter <file>",
	Short: "Filter .env entries by prefix, suffix, key list, or secret status",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]

		entries, err := envfile.Parse(path)
		if err != nil {
			return fmt.Errorf("parse error: %w", err)
		}

		prefix, _ := cmd.Flags().GetString("prefix")
		suffix, _ := cmd.Flags().GetString("suffix")
		keyList, _ := cmd.Flags().GetString("keys")
		secretsOnly, _ := cmd.Flags().GetBool("secrets-only")
		nonSecretsOnly, _ := cmd.Flags().GetBool("non-secrets")
		maskSecrets, _ := cmd.Flags().GetBool("mask")

		var keys []string
		if keyList != "" {
			for _, k := range strings.Split(keyList, ",") {
				if k = strings.TrimSpace(k); k != "" {
					keys = append(keys, k)
				}
			}
		}

		opts := envfile.FilterOptions{
			Prefix:         prefix,
			Suffix:         suffix,
			Keys:           keys,
			SecretsOnly:    secretsOnly,
			NonSecretsOnly: nonSecretsOnly,
		}

		result := envfile.Filter(entries, opts)
		if len(result) == 0 {
			fmt.Fprintln(os.Stderr, "no entries matched the filter criteria")
			return nil
		}

		for _, e := range result {
			v := e.Value
			if maskSecrets && envfile.IsSecret(e.Key) {
				v = "***"
			}
			fmt.Printf("%s=%s\n", e.Key, v)
		}
		return nil
	},
}

func init() {
	filterCmd.Flags().String("prefix", "", "filter keys by prefix")
	filterCmd.Flags().String("suffix", "", "filter keys by suffix")
	filterCmd.Flags().String("keys", "", "comma-separated list of keys to include")
	filterCmd.Flags().Bool("secrets-only", false, "include only secret entries")
	filterCmd.Flags().Bool("non-secrets", false, "include only non-secret entries")
	filterCmd.Flags().Bool("mask", false, "mask secret values in output")
	rootCmd.AddCommand(filterCmd)
}
