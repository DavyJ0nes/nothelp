package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

type TrainingSession struct {
	Day         string `json:"day"`
	Type        string `json:"type"`
	Duration    string `json:"duration"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

type TrainingWeek struct {
	WeekNumber   int               `json:"week_number"`
	StartDate    string            `json:"start_date"`
	EndDate      string            `json:"end_date"`
	Goal         string            `json:"goal"`
	Volume       string            `json:"volume"`
	TargetWeight string            `json:"target_weight"`
	Sessions     []TrainingSession `json:"sessions"`
}

type TrainingPlan struct {
	Weeks []TrainingWeek `json:"weeks"`
}

func TrainingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "training",
		Short: "Track Vätternrundan training plan",
		Long:  `View and log your 19-week training plan for Vätternrundan 315km`,
	}

	cmd.AddCommand(TrainingListCmd())
	cmd.AddCommand(TrainingLogCmd())

	return cmd
}

func TrainingListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list [week|today|tomorrow]",
		Short: "Show training schedule",
		Long:  "Show current week, specific week number (1-19), today's session, or tomorrow's session",
		RunE: func(_ *cobra.Command, args []string) error {
			if len(args) == 0 {
				return trainingListRun(0)
			}

			arg := args[0]
			if arg == "today" {
				return trainingShowDay(0)
			}
			if arg == "tomorrow" {
				return trainingShowDay(1)
			}

			// Try to parse as week number
			weekNum := 0
			_, _ = fmt.Sscanf(arg, "%d", &weekNum)
			return trainingListRun(weekNum)
		},
	}
	return cmd
}

func TrainingLogCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "log [day]",
		Short: "Mark today's or specific day's session as completed",
		Long:  "Mark session complete. Use day name (monday, tuesday, etc) or leave empty for today",
		RunE: func(_ *cobra.Command, args []string) error {
			day := ""
			if len(args) > 0 {
				day = args[0]
			}
			return trainingLogRun(day)
		},
	}
	return cmd
}

func trainingListRun(weekNum int) error {
	plan, err := loadTrainingPlan()
	if err != nil {
		return err
	}

	// Determine which week to show
	var targetWeek *TrainingWeek
	if weekNum > 0 {
		// Show specific week
		for i := range plan.Weeks {
			if plan.Weeks[i].WeekNumber == weekNum {
				targetWeek = &plan.Weeks[i]
				break
			}
		}
		if targetWeek == nil {
			return fmt.Errorf("week %d not found in plan", weekNum)
		}
	} else {
		// Show current week based on date
		today := time.Now()
		for i := range plan.Weeks {
			start, _ := time.Parse("2006-01-02", plan.Weeks[i].StartDate)
			end, _ := time.Parse("2006-01-02", plan.Weeks[i].EndDate)
			if (today.Equal(start) || today.After(start)) && (today.Equal(end) || today.Before(end)) {
				targetWeek = &plan.Weeks[i]
				break
			}
		}
		if targetWeek == nil {
			return fmt.Errorf("current week not found in plan")
		}
	}

	// Display the week
	fmt.Printf("\n🚴 WEEK %d: %s\n", targetWeek.WeekNumber, targetWeek.Goal)
	fmt.Printf("📅 %s → %s\n", targetWeek.StartDate, targetWeek.EndDate)
	fmt.Printf("⏱️  Volume: %s | 🎯 Target Weight: %s\n\n", targetWeek.Volume, targetWeek.TargetWeight)

	// Show sessions
	for _, session := range targetWeek.Sessions {
		status := "⬜"
		if session.Completed {
			status = "✅"
		}

		fmt.Printf("%s %s:\n", status, session.Day)
		fmt.Printf("   %s - %s\n", session.Type, session.Duration)
		if session.Description != "" {
			fmt.Printf("   %s\n", session.Description)
		}
		fmt.Println()
	}

	// Calculate progress
	completed := 0
	for _, s := range targetWeek.Sessions {
		if s.Completed {
			completed++
		}
	}
	fmt.Printf("Progress: %d/%d sessions completed\n\n", completed, len(targetWeek.Sessions))

	return nil
}

func trainingShowDay(daysOffset int) error {
	plan, err := loadTrainingPlan()
	if err != nil {
		return err
	}

	// Calculate target date
	targetDate := time.Now().AddDate(0, 0, daysOffset)
	targetDateStr := targetDate.Format("2006-01-02")
	targetDayName := targetDate.Format("Monday")

	// Find the week containing target date
	var targetWeek *TrainingWeek
	for i := range plan.Weeks {
		start, _ := time.Parse("2006-01-02", plan.Weeks[i].StartDate)
		end, _ := time.Parse("2006-01-02", plan.Weeks[i].EndDate)
		if (targetDate.Equal(start) || targetDate.After(start)) && (targetDate.Equal(end) || targetDate.Before(end)) {
			targetWeek = &plan.Weeks[i]
			break
		}
	}

	if targetWeek == nil {
		return fmt.Errorf("no training week found for %s", targetDateStr)
	}

	// Find session for target day
	var targetSession *TrainingSession
	for i := range targetWeek.Sessions {
		if targetWeek.Sessions[i].Day == targetDayName {
			targetSession = &targetWeek.Sessions[i]
			break
		}
	}

	if targetSession == nil {
		return fmt.Errorf("no session found for %s in week %d", targetDayName, targetWeek.WeekNumber)
	}

	// Display the session
	dayLabel := "TODAY"
	if daysOffset == 1 {
		dayLabel = "TOMORROW"
	}

	fmt.Printf("\n🚴 %s - %s (%s)\n", dayLabel, targetDayName, targetDateStr)
	fmt.Printf("Week %d: %s\n\n", targetWeek.WeekNumber, targetWeek.Goal)

	status := "⬜ Not completed"
	if targetSession.Completed {
		status = "✅ Completed"
	}

	fmt.Printf("%s\n", status)
	fmt.Printf("Type: %s\n", targetSession.Type)
	fmt.Printf("Duration: %s\n", targetSession.Duration)
	if targetSession.Description != "" {
		fmt.Printf("Details: %s\n", targetSession.Description)
	}
	fmt.Println()

	return nil
}

func trainingLogRun(day string) error {
	plan, err := loadTrainingPlan()
	if err != nil {
		return err
	}

	// Determine day to log
	targetDay := day
	if targetDay == "" {
		targetDay = time.Now().Format("Monday")
	}

	// Find current week
	today := time.Now()
	var currentWeek *TrainingWeek
	var weekIndex int
	for i := range plan.Weeks {
		start, _ := time.Parse("2006-01-02", plan.Weeks[i].StartDate)
		end, _ := time.Parse("2006-01-02", plan.Weeks[i].EndDate)
		if (today.Equal(start) || today.After(start)) && (today.Equal(end) || today.Before(end)) {
			currentWeek = &plan.Weeks[i]
			weekIndex = i
			break
		}
	}

	if currentWeek == nil {
		return fmt.Errorf("current week not found in plan")
	}

	// Find and mark session as completed
	found := false
	for i := range currentWeek.Sessions {
		if currentWeek.Sessions[i].Day == targetDay {
			currentWeek.Sessions[i].Completed = true
			found = true
			fmt.Printf("✅ Marked %s session complete: %s - %s\n",
				targetDay, currentWeek.Sessions[i].Type, currentWeek.Sessions[i].Duration)
			break
		}
	}

	if !found {
		return fmt.Errorf("no session found for %s in week %d", targetDay, currentWeek.WeekNumber)
	}

	// Update the plan in memory and save
	plan.Weeks[weekIndex] = *currentWeek
	if err := saveTrainingPlan(plan); err != nil {
		return err
	}

	// Show updated progress
	completed := 0
	for _, s := range currentWeek.Sessions {
		if s.Completed {
			completed++
		}
	}
	fmt.Printf("Week %d progress: %d/%d sessions completed\n", currentWeek.WeekNumber, completed, len(currentWeek.Sessions))

	return nil
}

func getTrainingPlanPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".habittracker", "training_plan.json")
}

func loadTrainingPlan() (TrainingPlan, error) {
	path := getTrainingPlanPath()
	plan := TrainingPlan{}

	file, err := os.ReadFile(path)
	if err != nil {
		return plan, fmt.Errorf("training plan not found. Create %s first", path)
	}

	if err := json.Unmarshal(file, &plan); err != nil {
		return plan, err
	}

	return plan, nil
}

func saveTrainingPlan(plan TrainingPlan) error {
	path := getTrainingPlanPath()

	file, err := json.MarshalIndent(plan, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, file, 0644)
}
