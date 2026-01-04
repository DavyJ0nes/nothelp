//nolint:gosec
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/davyj0nes/nothelp/internal/config"
	"github.com/davyj0nes/nothelp/internal/templates"
)

const dateFormat = "2006-01-02"

func openNoteFile(offset int, searchText, date string) error {
	conf, err := config.Parse()
	if err != nil {
		return err
	}

	filePath, err := getFilePath(conf, date)
	if err != nil {
		return err
	}

	lineNumber, err := templates.GetLineNumber(filePath, searchText)
	if err != nil {
		return err
	}

	return openInNvim(filePath, lineNumber+offset)
}

func createNewFile(conf config.Config, date, filePath string) error {
	body, err := templates.Parse(date)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(conf.DataLocation, 0o755); err != nil {
		return err
	}

	if err := os.WriteFile(filePath, body, 0o600); err != nil {
		return err
	}

	return nil
}

func getFilePath(conf config.Config, date string) (string, error) {
	filePath := conf.GetDataFilePath(date)
	if !exists(filePath) {
		if archiveFilePath, ok := fileInArchive(conf, date); ok {
			return archiveFilePath, nil
		}
		if err := createNewFile(conf, date, filePath); err != nil {
			return "", err
		}
	}
	return filePath, nil
}

func fileInArchive(conf config.Config, date string) (string, bool) {
	filePath := conf.GetArchiveFilePath(date)
	return filePath, exists(filePath)
}

func exists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

//nolint:gosec
func openInNvim(filePath string, lineNumber int) error {
	cmd := exec.Command("nvim", fmt.Sprintf("+%d", lineNumber), filePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func killApps(apps []string) error {
	for _, app := range apps {
		// Using pkill -i for case-insensitive matching
		if err := exec.Command("osascript", "-e", `tell application "`+app+`" to quit`).Run(); err != nil {
			return err
		}
	}
	return nil
}

func startApps(apps []string) error {
	for _, app := range apps {
		// Using pkill -i for case-insensitive matching
		if err := exec.Command("osascript", "-e", `tell application "`+app+`" to launch`).Run(); err != nil {
			return err
		}
	}
	return nil
}

func todayDate() string {
	return time.Now().Format(dateFormat)
}

func yesterdayDate() string {
	return time.Now().Add(-24 * time.Hour).Format(dateFormat)
}
