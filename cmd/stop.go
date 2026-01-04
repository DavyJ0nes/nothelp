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
	fmt.Println("🏁 Killing work apps...")
	err := killApps([]string{"Slack", "Arc"})
	if err != nil {
		return err
	}

	err = openNoteFile(0, "## 🏁 Shutdown", todayDate())
	if err != nil {
		return err
	}

	fmt.Println("\n--------------------------------------------------")
	fmt.Println("✅ Work context unloaded.")
	fmt.Println("💡 Remember: Your value is not tied to your output.")
	fmt.Println("🚴 Go cycle. 🍱 Go cook. Enjoy your life.")
	fmt.Println("--------------------------------------------------")

	return nil
}
