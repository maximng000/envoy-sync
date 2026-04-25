package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envoy-sync/internal/envfile"
)

var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Manage and compare environment profiles",
}

var profileListCmd = &cobra.Command{
	Use:   "list [dir]",
	Short: "List available profiles in a directory",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dir := "."
		if len(args) == 1 {
			dir = args[0]
		}
		profiles, err := envfile.ListProfiles(dir)
		if err != nil {
			return err
		}
		if len(profiles) == 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "No profiles found.")
			return nil
		}
		fmt.Fprintln(cmd.OutOrStdout(), strings.Join(profiles, "\n"))
		return nil
	},
}

var profileDiffCmd = &cobra.Command{
	Use:   "diff <dir> <profileA> <profileB>",
	Short: "Diff two environment profiles",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, nameA, nameB := args[0], args[1], args[2]

		a, err := envfile.LoadProfile(dir, nameA)
		if err != nil {
			return err
		}
		b, err := envfile.LoadProfile(dir, nameB)
		if err != nil {
			return err
		}

		diffs := envfile.DiffProfiles(a, b)
		if len(diffs) == 0 {
			fmt.Fprintf(cmd.OutOrStdout(), "Profiles %q and %q are identical.\n", nameA, nameB)
			return nil
		}

		for _, d := range diffs {
			switch d.Type {
			case "added":
				fmt.Fprintf(cmd.OutOrStdout(), "+ %s=%s\n", d.Key, d.NewValue)
			case "removed":
				fmt.Fprintf(cmd.OutOrStdout(), "- %s=%s\n", d.Key, d.OldValue)
			case "changed":
				fmt.Fprintf(cmd.OutOrStdout(), "~ %s: %s → %s\n", d.Key, d.OldValue, d.NewValue)
			}
		}
		os.Exit(1)
		return nil
	},
}

func init() {
	profileCmd.AddCommand(profileListCmd)
	profileCmd.AddCommand(profileDiffCmd)
	rootCmd.AddCommand(profileCmd)
}
