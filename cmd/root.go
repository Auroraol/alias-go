package cmd

import (
	"github.com/spf13/cobra"
)

// 根命令
var RootCmd = &cobra.Command{
	Use:   "als",
	Short: "A utility to manage aliases across shells",
	Long:  `A utility to manage aliases across shells with support for bash, zsh, fish, and PowerShell.`,
}
