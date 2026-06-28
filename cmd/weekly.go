package cmd

import (
	"github.com/spf13/cobra"
)

func WeeklyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "weekly",
		Short:   "open this week's review note",
		Aliases: []string{"week"},
		RunE: func(_ *cobra.Command, _ []string) error {
			return weeklyRun()
		},
	}
	return cmd
}

func weeklyRun() error {
	return openWeeklyNoteFile(0, "# Weekly Review", thisWeek())
}
