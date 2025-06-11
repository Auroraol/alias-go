package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

// 通用错误处理函数
func HandleError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		os.Exit(1)
	}
}

// 通用执行函数包装器
func ExecuteCommand(fn func() error) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		HandleError(fn())
	}
}
