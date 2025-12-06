package cmd

import (
	"github.com/spf13/cobra"
)

func StopCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "stop for the day",
		RunE: func(_ *cobra.Command, _ []string) error {
			return stopRun()
		},
	}
	return cmd
}

func stopRun() error {
	return openNoteFile(4, "## Evening Checklist", todayDate())
}
