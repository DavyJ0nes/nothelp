package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

type PressupEntry struct {
	Date      string `json:"date"`
	Morning   int    `json:"morning"`
	Evening   int    `json:"evening"`
	Total     int    `json:"total"`
	Completed bool   `json:"completed"`
}

type PressupData struct {
	Entries []PressupEntry `json:"entries"`
}

func PressupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pressup",
		Short: "Track pressup challenge progress",
		Long:  `Log and view your 50 pressups daily challenge (5x10 split 3 morning/2 evening)`,
	}

	cmd.AddCommand(PressupLogCmd())
	cmd.AddCommand(PressupListCmd())

	return cmd
}

func PressupLogCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "log [morning|evening] [rounds]",
		Short: "Log pressup session",
		Args:  cobra.ExactArgs(2),
		RunE: func(_ *cobra.Command, args []string) error {
			return pressupLogRun(args[0], args[1])
		},
	}
	return cmd
}

func PressupListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Show pressup challenge history and streak",
		RunE: func(_ *cobra.Command, _ []string) error {
			return pressupListRun()
		},
	}
	return cmd
}

func pressupLogRun(session, roundsStr string) error {
	rounds := 0
	_, _ = fmt.Sscanf(roundsStr, "%d", &rounds)

	data, err := loadPressupData()
	if err != nil {
		return err
	}

	today := time.Now().Format("2006-01-02")

	// Find or create today's entry
	var entry *PressupEntry
	for i := range data.Entries {
		if data.Entries[i].Date == today {
			entry = &data.Entries[i]
			break
		}
	}

	if entry == nil {
		newEntry := PressupEntry{Date: today}
		data.Entries = append(data.Entries, newEntry)
		entry = &data.Entries[len(data.Entries)-1]
	}

	// Update session
	switch session {
	case "morning":
		entry.Morning = rounds
	case "evening":
		entry.Evening = rounds
	}

	// Calculate total (10 pushups per round)
	entry.Total = (entry.Morning + entry.Evening) * 10
	entry.Completed = entry.Total >= 50

	if err := savePressupData(data); err != nil {
		return err
	}

	fmt.Printf("✓ Logged %d rounds (%s session)\n", rounds, session)
	fmt.Printf("Today's total: %d/%d pushups", entry.Total, 50)
	if entry.Completed {
		fmt.Printf(" 🎯 TARGET HIT!\n")
	} else {
		fmt.Printf("\n")
	}

	return nil
}

func pressupListRun() error {
	data, err := loadPressupData()
	if err != nil {
		return err
	}

	if len(data.Entries) == 0 {
		fmt.Println("No entries yet. Start logging with: pressup log morning 3")
		return nil
	}

	// Calculate streak
	streak := calculateStreak(data.Entries)
	totalDays := countCompletedDays(data.Entries)

	fmt.Println("\n🔥 PRESSUP CHALLENGE")
	fmt.Printf("Current Streak: %d days\n", streak)
	fmt.Printf("Total Completed: %d/30 days\n", totalDays)
	fmt.Println("\nRecent History:")
	fmt.Println("Date       | Morning | Evening | Total | Status")
	fmt.Println("-----------|---------|---------|-------|--------")

	// Show last 10 entries
	start := max(len(data.Entries)-10, 0)

	for i := len(data.Entries) - 1; i >= start; i-- {
		e := data.Entries[i]
		status := "❌"
		if e.Completed {
			status = "✅"
		}
		fmt.Printf("%s | %7d | %7d | %5d | %s\n",
			e.Date, e.Morning, e.Evening, e.Total, status)
	}
	fmt.Println()

	return nil
}

func getPressupDataPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".habittracker", "pressup.json")
}

func loadPressupData() (PressupData, error) {
	path := getPressupDataPath()
	data := PressupData{Entries: []PressupEntry{}}

	file, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return data, nil // Return empty if file doesn't exist
		}
		return data, err
	}

	if err := json.Unmarshal(file, &data); err != nil {
		return data, err
	}

	return data, nil
}

func savePressupData(data PressupData) error {
	path := getPressupDataPath()

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	file, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, file, 0644)
}

func calculateStreak(entries []PressupEntry) int {
	if len(entries) == 0 {
		return 0
	}

	streak := 0
	today := time.Now()

	// Walk backwards from today
	for i := 0; i <= len(entries); i++ {
		checkDate := today.AddDate(0, 0, -i).Format("2006-01-02")
		found := false

		for _, e := range entries {
			if e.Date == checkDate && e.Completed {
				streak++
				found = true
				break
			}
		}

		if !found {
			break
		}
	}

	return streak
}

func countCompletedDays(entries []PressupEntry) int {
	count := 0
	for _, e := range entries {
		if e.Completed {
			count++
		}
	}
	return count
}
