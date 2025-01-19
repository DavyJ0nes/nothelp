package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/davyj0nes/nothelp/internal/config"
	"github.com/davyj0nes/nothelp/internal/templates"
)

func exists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

func openInNvim(filePath string, lineNumber int) error {
	cmd := exec.Command("nvim", fmt.Sprintf("+%d", lineNumber), filePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func openNoteFile(offset int, searchText, date string) error {
	conf, err := config.Parse()
	if err != nil {
		return err
	}

	filePath := conf.DataLocation + "/" + date + ".md"

	if !exists(filePath) {
		body, err := templates.Parse(date)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(conf.DataLocation, 0755); err != nil {
			return err
		}
		if err := os.WriteFile(filePath, body, 0644); err != nil {
			return err
		}
	}

	lineNumber, err := templates.GetLineNumber(filePath, searchText)
	if err != nil {
		return err
	}

	return openInNvim(filePath, lineNumber+offset)
}
