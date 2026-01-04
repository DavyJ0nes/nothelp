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
	err := killApps([]string{"Steam", "Firefox", "Chrome", "Music", "Preview"})
	if err != nil {
		return err
	}

	err = startApps([]string{"Slack", "Arc"})
	if err != nil {
		return err
	}

	return openNoteFile(0, "## 🚀 Startup", todayDate())
}
