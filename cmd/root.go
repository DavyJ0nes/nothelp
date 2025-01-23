package cmd

import (
	"github.com/spf13/cobra"
)

func RootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nothelp",
		Short: "note helper",
		Long:  `Simple application to help you take notes.`,
	}
	cmd.AddCommand(
		StartCmd(),
		StopCmd(),
		TodayCmd(),
		YesterdayCmd(),
	)
	return cmd
}
