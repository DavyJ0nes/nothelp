package cmd

import (
	"github.com/spf13/cobra"
)

func NoteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "notes",
		Short:   "open notes",
		Aliases: []string{"notes"},
		RunE: func(_ *cobra.Command, _ []string) error {
			return noteRun()
		},
	}
	return cmd
}

func noteRun() error {
	return openNoteFile(-1, "## 📝 The Log", todayDate())
}
