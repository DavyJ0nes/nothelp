package cmd

import (
	"time"

	"github.com/spf13/cobra"
)

func StopCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "stop for the day",
		RunE: func(cmd *cobra.Command, args []string) error {
			return stopRun()
		},
	}
	return cmd
}

func stopRun() error {
	date := time.Now().Format("2006-01-02")
	return openNoteFile(3, "## Evening Checklist", date)
}
