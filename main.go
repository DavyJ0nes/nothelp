package main

import (
	"fmt"
	"os"

	"github.com/davyj0nes/nothelp/cmd"
)

func main() {
	if err := cmd.RootCmd().Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "command failed to run: %s\n", err)
		os.Exit(1)
	}
}
