package cmd

import (
	"github.com/spf13/cobra"
)

func TodayCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "today",
		Short:   "open todays note",
		Aliases: []string{"inbox"},
		RunE: func(_ *cobra.Command, _ []string) error {
			return todayRun()
		},
	}
	return cmd
}

func todayRun() error {
	return openNoteFile(-1, "## Focus for Today", todayDate())
}
