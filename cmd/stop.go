package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func StopCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "stop for the day",
		RunE: func(_ *cobra.Command, _ []string) error {
			return stopRun()
		},
	}
	return cmd
}

func stopRun() error {
	if err := openDailyNoteFile(0, "## 🏁 Shutdown", todayDate()); err != nil {
		return err
	}

	fmt.Println("\n--------------------------------------------------")
	fmt.Println("✅ Work context unloaded.")
	fmt.Println("💡 Remember: Your value is not tied to your output.")
	fmt.Println("🚴 Go cycle. 🍱 Go cook. Enjoy your life.")
	fmt.Println("--------------------------------------------------")

	return nil
}
