package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-sync/internal/envfile"
)

var promoteCmd = &cobra.Command{
	Use:   "promote <src> <dst>",
	Short: "Promote env variables from one environment file to another",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		overwrite, _ := cmd.Flags().GetBool("overwrite")
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		maskSecrets, _ := cmd.Flags().GetBool("mask-secrets")

		srcEntries, err := envfile.Parse(args[0])
		if err != nil {
			return fmt.Errorf("reading src: %w", err)
		}
		dstEntries, err := envfile.Parse(args[1])
		if err != nil {
			return fmt.Errorf("reading dst: %w", err)
		}

		opts := envfile.PromoteOptions{
			Overwrite:   overwrite,
			DryRun:      dryRun,
			MaskSecrets: maskSecrets,
		}

		out, results, err := envfile.Promote(srcEntries, dstEntries, opts)
		if err != nil {
			return err
		}

		for _, r := range results {
			switch r.Action {
			case "added":
				fmt.Fprintf(cmd.OutOrStdout(), "+ %s=%s\n", r.Key, r.NewValue)
			case "updated":
				fmt.Fprintf(cmd.OutOrStdout(), "~ %s: %s -> %s\n", r.Key, r.OldValue, r.NewValue)
			case "skipped":
				fmt.Fprintf(cmd.OutOrStdout(), "! %s skipped (use --overwrite to update)\n", r.Key)
			case "unchanged":
				fmt.Fprintf(cmd.OutOrStdout(), "= %s unchanged\n", r.Key)
			}
		}

		if dryRun {
			fmt.Fprintln(cmd.OutOrStdout(), "[dry-run] no changes written")
			return nil
		}

		f, err := os.Create(args[1])
		if err != nil {
			return fmt.Errorf("writing dst: %w", err)
		}
		defer f.Close()
		for _, e := range out {
			fmt.Fprintf(f, "%s=%s\n", e.Key, e.Value)
		}
		return nil
	},
}

func init() {
	promoteCmd.Flags().Bool("overwrite", false, "overwrite existing keys in dst")
	promoteCmd.Flags().Bool("dry-run", false, "preview changes without writing")
	promoteCmd.Flags().Bool("mask-secrets", true, "mask secret values in output")
	rootCmd.AddCommand(promoteCmd)
}
