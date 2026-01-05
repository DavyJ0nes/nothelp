package cmd

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/davyj0nes/nothelp/internal/config"
	"github.com/spf13/cobra"
)

const (
	timeAuditHeader = "## ⏳ The Time Audit"
	logHeader       = "## 📝 The Log"
)

type auditEntry struct {
	Timestamp time.Time
	Activity  string
	Minutes   int
	Category  string
}

func AnalysisCmd() *cobra.Command {
	var dateFlag string

	cmd := &cobra.Command{
		Use:   "analysis [date]",
		Short: "analyze a daily time audit for focus scoring",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			date := dateFlag
			if len(args) == 1 {
				if dateFlag != "" {
					return errors.New("provide the date either as an argument or --date, not both")
				}
				date = args[0]
			}
			if date == "" {
				date = todayDate()
			}
			if _, err := time.Parse(dateFormat, date); err != nil {
				return fmt.Errorf("invalid date %q, expected YYYY-MM-DD", date)
			}

			return analysisRun(cmd, date)
		},
	}

	cmd.Flags().StringVarP(&dateFlag, "date", "d", "", "date to analyze (YYYY-MM-DD)")

	return cmd
}

func analysisRun(cmd *cobra.Command, date string) error {
	conf, err := config.Parse()
	if err != nil {
		return err
	}

	body, found, err := readNoteForDate(conf, date)
	if err != nil {
		return err
	}
	if !found {
		return fmt.Errorf("no note found for %s", date)
	}

	lines := strings.Split(body, "\n")
	section, err := extractTimeAuditSection(lines)
	if err != nil {
		return err
	}

	entries, err := parseTimeAuditEntries(section)
	if err != nil {
		return err
	}

	entries = applyDurations(entries)
	summary := buildFocusSummary(entries, len(entries))

	fmt.Fprintln(cmd.OutOrStdout(), summary)
	return nil
}

func extractTimeAuditSection(lines []string) ([]string, error) {
	start := sectionHeaderIndex(lines, timeAuditHeader)
	if start == -1 {
		return nil, fmt.Errorf("time audit section not found")
	}

	end := sectionHeaderIndex(lines, logHeader)
	if end == -1 {
		end = sectionBoundaryIndex(lines, start+1)
		if end == -1 {
			return nil, fmt.Errorf("time audit section boundary not found")
		}
	}

	if end <= start+1 {
		return nil, fmt.Errorf("time audit section is empty")
	}

	return lines[start+1 : end], nil
}

func parseTimeAuditEntries(lines []string) ([]auditEntry, error) {
	entries := make([]auditEntry, 0, len(lines))

	for _, line := range lines {
		parts := parseMarkdownRow(line)
		if len(parts) < 2 {
			continue
		}

		timeStr := parts[0]
		activity := parts[1]
		if strings.EqualFold(timeStr, "time") || strings.EqualFold(activity, "activity") {
			continue
		}
		if isDividerCell(timeStr) || isPlaceholderTime(timeStr) {
			continue
		}

		parsed, err := time.ParseInLocation("15:04", timeStr, time.Local)
		if err != nil {
			return nil, fmt.Errorf("invalid time %q in time audit", timeStr)
		}

		entries = append(entries, auditEntry{
			Timestamp: parsed,
			Activity:  activity,
		})
	}

	if len(entries) == 0 {
		return nil, fmt.Errorf("no time audit entries found")
	}

	return entries, nil
}

func parseMarkdownRow(line string) []string {
	trimmed := strings.TrimSpace(line)
	if !strings.HasPrefix(trimmed, "|") {
		return nil
	}
	trimmed = strings.Trim(trimmed, "|")
	parts := strings.Split(trimmed, "|")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

func isDividerCell(value string) bool {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return true
	}
	for _, r := range trimmed {
		if r != ':' && r != '-' {
			return false
		}
	}
	return true
}

func isPlaceholderTime(value string) bool {
	lower := strings.ToLower(strings.TrimSpace(value))
	return strings.Contains(lower, "xx") || lower == "x"
}

func applyDurations(entries []auditEntry) []auditEntry {
	if len(entries) < 2 {
		return entries
	}

	for i := 0; i < len(entries)-1; i++ {
		diff := entries[i+1].Timestamp.Sub(entries[i].Timestamp)
		if diff < 0 {
			diff += 24 * time.Hour
		}
		entries[i].Minutes = int(diff.Minutes())
	}

	return entries
}

func buildFocusSummary(entries []auditEntry, entryCount int) string {
	gainKeywords := []string{
		"grid",
		"architecture",
		"elixir",
		"vättern",
		"training",
		"meditation",
		"logic",
	}
	lossKeywords := []string{
		"config",
		"nvim",
		"youtube",
		"slack",
		"browsing",
		"steam",
	}

	var gainMinutes int
	var lossMinutes int
	var lossOverLimit bool

	for _, entry := range entries {
		if entry.Minutes == 0 {
			continue
		}

		category := categorizeActivity(entry.Activity, gainKeywords, lossKeywords)
		switch category {
		case "gain":
			gainMinutes += entry.Minutes
		case "loss":
			lossMinutes += entry.Minutes
			if entry.Minutes > 30 {
				lossOverLimit = true
			}
		}
	}

	total := gainMinutes + lossMinutes
	focusScore := 0.0
	if total > 0 {
		focusScore = (float64(gainMinutes) / float64(total)) * 100
	}

	base := fmt.Sprintf(
		"Entries: %d | Gain: %d min | Loss: %d min | Focus Score: %.1f%%",
		entryCount,
		gainMinutes,
		lossMinutes,
		focusScore,
	)
	if lossOverLimit {
		return base + " | ⚠️"
	}
	return base
}

func categorizeActivity(activity string, gainKeywords, lossKeywords []string) string {
	lower := strings.ToLower(activity)
	for _, keyword := range gainKeywords {
		if strings.Contains(lower, keyword) {
			return "gain"
		}
	}
	for _, keyword := range lossKeywords {
		if strings.Contains(lower, keyword) {
			return "loss"
		}
	}
	return "neutral"
}
