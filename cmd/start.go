package cmd

import (
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
	return openNoteFile(4, "## Morning Checklist", todayDate())
}
