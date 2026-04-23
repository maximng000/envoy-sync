package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-sync/internal/envfile"
)

var encryptCmd = &cobra.Command{
	Use:   "encrypt <file>",
	Short: "Encrypt secret values in a .env file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath := args[0]
		key, _ := cmd.Flags().GetString("key")
		output, _ := cmd.Flags().GetString("output")
		decryptMode, _ := cmd.Flags().GetBool("decrypt")

		if len(key) != 16 && len(key) != 24 && len(key) != 32 {
			return fmt.Errorf("key must be 16, 24, or 32 characters long")
		}

		entries, err := envfile.Parse(filePath)
		if err != nil {
			return fmt.Errorf("failed to parse %s: %w", filePath, err)
		}

		var processed []envfile.Entry
		if decryptMode {
			processed, err = envfile.DecryptSecrets(entries, key)
			if err != nil {
				return fmt.Errorf("decryption failed: %w", err)
			}
		} else {
			processed, err = envfile.EncryptSecrets(entries, key)
			if err != nil {
				return fmt.Errorf("encryption failed: %w", err)
			}
		}

		dest := filePath
		if output != "" {
			dest = output
		}

		result, err := envfile.Export(processed, "dotenv", false)
		if err != nil {
			return fmt.Errorf("export failed: %w", err)
		}
		if err := os.WriteFile(dest, []byte(result), 0600); err != nil {
			return fmt.Errorf("failed to write %s: %w", dest, err)
		}

		action := "Encrypted"
		if decryptMode {
			action = "Decrypted"
		}
		fmt.Fprintf(cmd.OutOrStdout(), "%s secrets in %s -> %s\n", action, filePath, dest)
		return nil
	},
}

func init() {
	encryptCmd.Flags().StringP("key", "k", "", "AES encryption key (16, 24, or 32 chars) (required)")
	encryptCmd.Flags().StringP("output", "o", "", "Output file path (defaults to input file)")
	encryptCmd.Flags().Bool("decrypt", false, "Decrypt secrets instead of encrypting")
	_ = encryptCmd.MarkFlagRequired("key")
	rootCmd.AddCommand(encryptCmd)
}
