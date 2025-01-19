package main

import (
	"fmt"
	"os"

	"github.com/davyj0nes/nothelp/cmd"
)

func main() {
	if err := cmd.RootCmd().Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "nothelp failed: %s\n", err)
		os.Exit(1)
	}
}
