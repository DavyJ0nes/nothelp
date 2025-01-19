package cmd

import (
	"time"

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
	date := time.Now().Format("2006-01-02")
	return openNoteFile(
		-1,
		"______________________________________________________________________",
		date,
	)
}
