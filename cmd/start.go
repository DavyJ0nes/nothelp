package cmd

import (
	"time"

	"github.com/spf13/cobra"
)

func StartCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "start the day",
		RunE: func(_ *cobra.Command, _ []string) error {
			return startRun()
		},
	}
	return cmd
}

func startRun() error {
	date := time.Now().Format("2006-01-02")
	return openNoteFile(4, "## Morning Checklist", date)
}
