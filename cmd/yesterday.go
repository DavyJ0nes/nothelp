package cmd

import (
	"time"

	"github.com/spf13/cobra"
)

func YesterdayCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "yesterday",
		Short: "open yesterdays note",
		RunE: func(cmd *cobra.Command, args []string) error {
			return yesterdayRun()
		},
	}
	return cmd
}

func yesterdayRun() error {
	date := time.Now().Add(24 * time.Hour).Format("2006-01-02")
	return openNoteFile(
		-1,
		"-------------------------------------------------------------------------------",
		date,
	)
}
