package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/davyj0nes/nothelp/internal/config"
	"github.com/spf13/cobra"
)

func LogCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "log <activity...>",
		Short: "Log a timestamped activity to the daily audit",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			activity := strings.Join(args, " ")
			return logRun(activity)
		},
	}
}

func logRun(activity string) error {
	currentTime := time.Now().Format("15:04")

	// 1. Get configuration and file path
	conf, err := config.Parse()
	if err != nil {
		return fmt.Errorf("config error: %w", err)
	}

	filePath, err := getFilePath(conf, todayDate())
	if err != nil {
		return fmt.Errorf("file path error: %w", err)
	}

	// 2. Read the existing note
	input, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("could not read note at %s: %w", filePath, err)
	}

	lines := strings.Split(string(input), "\n")
	newRow := fmt.Sprintf("| %s | %s | — |", currentTime, activity)

	updated, err := insertTimeAuditRow(lines, newRow)
	if err != nil {
		return fmt.Errorf("could not find the injection point (Time Audit section) in %s: %w", filePath, err)
	}

	// 4. Save the file
	err = os.WriteFile(filePath, []byte(strings.Join(updated, "\n")), 0o600)
	if err != nil {
		return fmt.Errorf("failed to save note: %w", err)
	}

	fmt.Printf("⚔️  Log recorded at %s: %s\n", currentTime, activity)
	return nil
}

func insertTimeAuditRow(lines []string, row string) ([]string, error) {
	sectionIndex := sectionHeaderIndex(lines, "## ⏳ The Time Audit")
	if sectionIndex == -1 {
		return nil, fmt.Errorf("time audit section not found")
	}

	boundaryIndex := sectionBoundaryIndex(lines, sectionIndex+1)
	if boundaryIndex == -1 {
		return nil, fmt.Errorf("time audit boundary not found")
	}

	// Clean up trailing whitespace/newlines before the table ends.
	for boundaryIndex > sectionIndex+1 && strings.TrimSpace(lines[boundaryIndex-1]) == "" {
		boundaryIndex--
	}

	updated := make([]string, 0, len(lines)+2)
	updated = append(updated, lines[:boundaryIndex]...)
	updated = append(updated, row, "")
	updated = append(updated, lines[boundaryIndex:]...)
	return updated, nil
}

func sectionHeaderIndex(lines []string, header string) int {
	for i, line := range lines {
		if strings.TrimSpace(line) == header {
			return i
		}
	}
	return -1
}

func sectionBoundaryIndex(lines []string, start int) int {
	for i := start; i < len(lines); i++ {
		line := lines[i]
		if strings.HasPrefix(line, "## ") || strings.HasPrefix(line, "___") {
			return i
		}
	}
	return -1
}
