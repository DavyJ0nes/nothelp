package cmd

import (
	"github.com/spf13/cobra"
)

func YesterdayCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "yesterday",
		Short: "open yesterdays note",
		RunE: func(_ *cobra.Command, _ []string) error {
			return yesterdayRun()
		},
	}
	return cmd
}

func yesterdayRun() error {
	return openDailyNoteFile(-1, "## 📝 Log", yesterdayDate())
}
