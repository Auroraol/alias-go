package main

import (
	cmdlib "alias-go/cmd"
	"alias-go/cmd/alias"
	"alias-go/cmd/cron"
	"fmt"
	"os"
)

func main() {
	alias.InitAlias()
	cron.InitCron()
	if err := cmdlib.RootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
