package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envoy-sync/internal/envfile"
)

var patchCmd = &cobra.Command{
	Use:   "patch <file>",
	Short: "Apply set/delete/rename operations to an env file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath := args[0]

		setOps, _ := cmd.Flags().GetStringArray("set")
		delOps, _ := cmd.Flags().GetStringArray("delete")
		renOps, _ := cmd.Flags().GetStringArray("rename")
		inPlace, _ := cmd.Flags().GetBool("in-place")

		entries, err := envfile.Parse(filePath)
		if err != nil {
			return fmt.Errorf("failed to parse %s: %w", filePath, err)
		}

		var ops []envfile.PatchOp

		for _, s := range setOps {
			parts := strings.SplitN(s, "=", 2)
			if len(parts) != 2 {
				return fmt.Errorf("invalid --set value %q: expected KEY=VALUE", s)
			}
			ops = append(ops, envfile.PatchOp{Op: "set", Key: parts[0], Value: parts[1]})
		}

		for _, k := range delOps {
			ops = append(ops, envfile.PatchOp{Op: "delete", Key: k})
		}

		for _, r := range renOps {
			parts := strings.SplitN(r, ":", 2)
			if len(parts) != 2 {
				return fmt.Errorf("invalid --rename value %q: expected OLD:NEW", r)
			}
			ops = append(ops, envfile.PatchOp{Op: "rename", Key: parts[0], NewKey: parts[1]})
		}

		updated, results, err := envfile.Patch(entries, ops)
		if err != nil {
			return err
		}

		for _, r := range results {
			if !r.Applied {
				fmt.Fprintf(cmd.ErrOrStderr(), "warn: op %s on %q not applied: %s\n", r.Op, r.Key, r.Reason)
			}
		}

		exported, err := envfile.Export(updated, "dotenv", false)
		if err != nil {
			return fmt.Errorf("failed to export: %w", err)
		}

		if inPlace {
			if err := os.WriteFile(filePath, []byte(exported), 0644); err != nil {
				return fmt.Errorf("failed to write %s: %w", filePath, err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "patched %s\n", filePath)
		} else {
			fmt.Fprint(cmd.OutOrStdout(), exported)
		}

		return nil
	},
}

func init() {
	patchCmd.Flags().StringArray("set", nil, "Set a key: KEY=VALUE (repeatable)")
	patchCmd.Flags().StringArray("delete", nil, "Delete a key (repeatable)")
	patchCmd.Flags().StringArray("rename", nil, "Rename a key: OLD:NEW (repeatable)")
	patchCmd.Flags().Bool("in-place", false, "Write changes back to the source file")
	rootCmd.AddCommand(patchCmd)
}
