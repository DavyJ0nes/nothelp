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

func getDailyFilePath(conf config.Config, date string) (string, error) {
	filePath := conf.GetDataFilePath(date)
	if !exists(filePath) {
		if archiveFilePath, ok := fileInArchive(conf, date); ok {
			return archiveFilePath, nil
		}
		if err := createNewDailyFile(conf, date, filePath); err != nil {
			return "", err
		}
	}
	return filePath, nil
}

func openDailyNoteFile(offset int, searchText, date string) error {
	conf, err := config.Parse()
	if err != nil {
		return err
	}

	filePath, err := getDailyFilePath(conf, date)
	if err != nil {
		return err
	}

	lineNumber, err := templates.GetLineNumber(filePath, searchText)
	if err != nil {
		return err
	}

	return openInNvim(filePath, lineNumber+offset)
}

func openWeeklyNoteFile(offset int, searchText, week string) error {
	conf, err := config.Parse()
	if err != nil {
		return err
	}

	filePath, err := getWeeklyFilePath(conf, week)
	if err != nil {
		return err
	}

	lineNumber, err := templates.GetLineNumber(filePath, searchText)
	if err != nil {
		return err
	}

	return openInNvim(filePath, lineNumber+offset)
}

func getWeeklyFilePath(conf config.Config, week string) (string, error) {
	filePath := conf.GetWeeklyFilePath(week)
	if !exists(filePath) {
		if err := createNewWeeklyFile(conf, week, filePath); err != nil {
			return "", err
		}
	}
	return filePath, nil
}

func createNewWeeklyFile(conf config.Config, week, filePath string) error {
	body, err := templates.ParseWeekly(week)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(conf.WeeklyLocation, 0o755); err != nil {
		return err
	}

	if err := os.WriteFile(filePath, body, 0o600); err != nil {
		return err
	}

	return nil
}

func createNewDailyFile(conf config.Config, date, filePath string) error {
	body, err := templates.ParseDaily(date)
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

func todayDate() string {
	return time.Now().Format(dateFormat)
}

func yesterdayDate() string {
	return time.Now().Add(-24 * time.Hour).Format(dateFormat)
}

func thisWeek() string {
	year, week := time.Now().ISOWeek()
	return fmt.Sprintf("%d-W%02d", year, week)
}
