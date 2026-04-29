package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"envoy-sync/internal/envfile"
)

var archiveCmd = &cobra.Command{
	Use:   "archive",
	Short: "Manage versioned archives of .env files",
}

var archiveSaveCmd = &cobra.Command{
	Use:   "save <env-file> <archive-file>",
	Short: "Append current env file state to an archive",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		label, _ := cmd.Flags().GetString("label")
		entries, err := envfile.Parse(args[0])
		if err != nil {
			return fmt.Errorf("parse: %w", err)
		}
		if err := envfile.AddToArchive(args[1], entries, label); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Saved snapshot to %s\n", args[1])
		return nil
	},
}

var archiveDiffCmd = &cobra.Command{
	Use:   "diff <archive-file> <from-index> <to-index>",
	Short: "Diff two versions stored in an archive",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		fromIdx, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid from-index: %w", err)
		}
		toIdx, err := strconv.Atoi(args[2])
		if err != nil {
			return fmt.Errorf("invalid to-index: %w", err)
		}
		arch, err := envfile.LoadArchive(args[0])
		if err != nil {
			return err
		}
		diffs, err := envfile.DiffArchiveVersions(arch, fromIdx, toIdx)
		if err != nil {
			return err
		}
		if len(diffs) == 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "No differences.")
			return nil
		}
		for _, d := range diffs {
			fmt.Fprintln(cmd.OutOrStdout(), d.String())
		}
		return nil
	},
}

var archiveListCmd = &cobra.Command{
	Use:   "list <archive-file>",
	Short: "List all versions in an archive",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		arch, err := envfile.LoadArchive(args[0])
		if err != nil {
			return err
		}
		if len(arch.Versions) == 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "No versions found.")
			return nil
		}
		for i, v := range arch.Versions {
			label := v.Label
			if label == "" {
				label = "(no label)"
			}
			fmt.Fprintf(cmd.OutOrStdout(), "[%d] %s  %s  keys=%d\n",
				i, v.Timestamp.Format("2006-01-02T15:04:05Z"), label, len(v.Entries))
		}
		return nil
	},
}

func init() {
	archiveSaveCmd.Flags().String("label", "", "Optional label for this snapshot")
	archiveCmd.AddCommand(archiveSaveCmd, archiveDiffCmd, archiveListCmd)
	rootCmd.AddCommand(archiveCmd)
	_ = os.Stderr
}
