package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"envoy-sync/internal/envfile"
)

var watchInterval int

var watchCmd = &cobra.Command{
	Use:   "watch <file>",
	Short: "Watch a .env file for changes and print a diff on each update",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]
		interval := time.Duration(watchInterval) * time.Millisecond

		current, err := envfile.Parse(path)
		if err != nil {
			return fmt.Errorf("watch: failed to parse %q: %w", path, err)
		}

		done := make(chan struct{})
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			<-sig
			close(done)
		}()

		ch, err := envfile.Watch(path, interval, done)
		if err != nil {
			return fmt.Errorf("watch: %w", err)
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Watching %s (interval: %s)...\n", path, interval)

		for event := range ch {
			if event.Err != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), "error: %v\n", event.Err)
				continue
			}
			if !event.Changed {
				continue
			}

			updated, err := envfile.Parse(path)
			if err != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), "parse error: %v\n", err)
				continue
			}

			results := envfile.Diff(current, updated)
			if len(results) == 0 {
				current = updated
				continue
			}

			fmt.Fprintf(cmd.OutOrStdout(), "\n[%s] Changes detected in %s:\n",
				time.Now().Format("15:04:05"), path)
			for _, d := range results {
				fmt.Fprintln(cmd.OutOrStdout(), d)
			}
			current = updated
		}

		fmt.Fprintln(cmd.OutOrStdout(), "\nWatch stopped.")
		return nil
	},
}

func init() {
	watchCmd.Flags().IntVar(&watchInterval, "interval", 1000, "Poll interval in milliseconds")
	rootCmd.AddCommand(watchCmd)
}
