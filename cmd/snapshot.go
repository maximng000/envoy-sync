package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-sync/internal/envfile"
)

var snapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Short: "Manage env file snapshots",
}

var snapshotTakeCmd = &cobra.Command{
	Use:   "take <env-file> <snapshot-file>",
	Short: "Take a snapshot of an env file",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		envPath, snapPath := args[0], args[1]

		entries, err := envfile.Parse(envPath)
		if err != nil {
			return fmt.Errorf("parse env file: %w", err)
		}

		snap := envfile.TakeSnapshot(envPath, entries)
		if err := envfile.SaveSnapshot(snapPath, snap); err != nil {
			return err
		}

		fmt.Fprintf(os.Stdout, "Snapshot saved to %s (%d keys)\n", snapPath, len(entries))
		return nil
	},
}

var snapshotDiffCmd = &cobra.Command{
	Use:   "diff <snapshot-file> <env-file>",
	Short: "Diff a snapshot against a current env file",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		snapPath, envPath := args[0], args[1]

		snap, err := envfile.LoadSnapshot(snapPath)
		if err != nil {
			return fmt.Errorf("load snapshot: %w", err)
		}

		current, err := envfile.Parse(envPath)
		if err != nil {
			return fmt.Errorf("parse env file: %w", err)
		}

		diffs := envfile.DiffSnapshot(snap, current)
		if len(diffs) == 0 {
			fmt.Fprintln(os.Stdout, "No differences from snapshot.")
			return nil
		}

		fmt.Fprintf(os.Stdout, "Snapshot: %s @ %s\n", snap.Source, snap.Timestamp.Format("2006-01-02 15:04:05 UTC"))
		for _, d := range diffs {
			fmt.Fprintln(os.Stdout, d.String())
		}
		return nil
	},
}

func init() {
	snapshotCmd.AddCommand(snapshotTakeCmd)
	snapshotCmd.AddCommand(snapshotDiffCmd)
	rootCmd.AddCommand(snapshotCmd)
}
