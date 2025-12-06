package cmd

import (
	"os"
	"path/filepath"

	"github.com/davyj0nes/nothelp/internal/config"
	"github.com/spf13/cobra"
)

func ArchiveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "archive",
		Short: "archive notes",
		RunE: func(_ *cobra.Command, _ []string) error {
			return archiveRun()
		},
	}
	return cmd
}

func archiveRun() error {
	conf, err := config.Parse()
	if err != nil {
		return err
	}

	entries, err := os.ReadDir(conf.DataLocation)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(conf.ArchiveLocation, 0o755); err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		src := filepath.Join(conf.DataLocation, entry.Name())
		dest := filepath.Join(conf.ArchiveLocation, entry.Name())

		if err := os.Rename(src, dest); err != nil {
			return err
		}
	}

	return nil
}
