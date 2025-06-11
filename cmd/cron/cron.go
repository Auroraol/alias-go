package cron

import (
	cmdlib "alias-go/cmd"
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// 常量定义
const (
	separatorLine  = "------------------------------------------------------------"
	cronTimeFormat = "15:04:05"
)

func InitCron() {
	// 注册主命令
	cmdlib.RootCmd.AddCommand(cronCmd)

	// 注册子命令 - 使用切片批量添加，更简洁
	subCommands := []*cobra.Command{
		cronListCmd,
		cronAddCmd,
		cronEditCmd,
		cronMonitorCmd,
		cronDeleteCmd,
	}

	for _, subCmd := range subCommands {
		cronCmd.AddCommand(subCmd)
	}
}

// 主 cron 命令
var cronCmd = &cobra.Command{
	Use:   "cron",
	Short: "管理 cron 任务",
	Long:  `提供查看、添加、编辑和监控 cron 任务的功能`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// 查看 cron 任务
var cronListCmd = &cobra.Command{
	Use:   "list",
	Short: "查看当前用户的 cron 任务",
	Long:  `显示当前用户的所有 cron 任务`,
	Run:   cmdlib.ExecuteCommand(listCronJobs),
}

// 添加 cron 任务
var cronAddCmd = &cobra.Command{
	Use:   "add [时间表达式] [命令]",
	Short: "添加新的 cron 任务",
	Long: `添加新的 cron 任务。时间表达式格式: 分 时 日 月 周
示例: als cron add "0 2 * * *" "backup.sh"`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		schedule, command := args[0], args[1]

		if err := addCronJob(schedule, command); err != nil {
			cmdlib.HandleError(err)
		}
		fmt.Printf("成功添加 cron 任务: %s %s\n", schedule, command)
	},
}

// 编辑 cron 任务
var cronEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "编辑 cron 任务",
	Long:  `打开系统默认编辑器编辑 crontab`,
	Run:   cmdlib.ExecuteCommand(editCronJobs),
}

// 监控 cron 任务
var cronMonitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "监控最近的 cron 任务执行情况",
	Long:  `实时监控 cron 任务的执行情况和日志`,
	Run:   cmdlib.ExecuteCommand(monitorCronJobs),
}

// 删除 cron 任务
var cronDeleteCmd = &cobra.Command{
	Use:   "delete [行号]",
	Short: "删除指定的 cron 任务",
	Long:  `根据行号删除指定的 cron 任务。先使用 'als cron list' 查看行号`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		lineNum, err := strconv.Atoi(args[0])
		if err != nil {
			cmdlib.HandleError(fmt.Errorf("无效的行号 '%s'", args[0]))
		}

		if err := deleteCronJob(lineNum); err != nil {
			cmdlib.HandleError(err)
		}
		fmt.Printf("成功删除第 %d 行的 cron 任务\n", lineNum)
	},
}

// =============================================================================
// 核心功能实现
// =============================================================================

// listCronJobs 列出当前用户的 cron 任务
func listCronJobs() error {
	if runtime.GOOS == "windows" {
		return listWindowsScheduledTasks()
	}
	return listUnixCronJobs()
}

// listUnixCronJobs 列出Unix系统的cron任务
func listUnixCronJobs() error {
	cmd := exec.Command("crontab", "-l")
	output, err := cmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok && exitError.ExitCode() == 1 {
			fmt.Println("当前用户没有 cron 任务")
			return nil
		}
		return fmt.Errorf("获取 cron 任务失败: %v", err)
	}

	return displayCronJobs(string(output))
}

// displayCronJobs 显示cron任务列表
func displayCronJobs(output string) error {
	lines := strings.Split(output, "\n")
	fmt.Println("当前 cron 任务:")
	fmt.Println(separatorLine)

	jobCount := 0
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "#") {
			fmt.Printf("%d: %s (注释)\n", i+1, line)
		} else {
			jobCount++
			fmt.Printf("%d: %s\n", i+1, line)
		}
	}

	if jobCount == 0 {
		fmt.Println("没有活跃的 cron 任务")
	} else {
		fmt.Printf("\n总共 %d 个活跃任务\n", jobCount)
	}

	return nil
}

// addCronJob 添加新的 cron 任务
func addCronJob(schedule, command string) error {
	if runtime.GOOS == "windows" {
		return addWindowsScheduledTask(schedule, command)
	}
	return addUnixCronJob(schedule, command)
}

// addUnixCronJob 添加Unix系统的cron任务
func addUnixCronJob(schedule, command string) error {
	// 验证 cron 表达式
	if !isValidCronExpression(schedule) {
		return fmt.Errorf("无效的 cron 表达式: %s", schedule)
	}

	// 获取现有的 crontab
	currentCrontab, err := getCurrentCrontab()
	if err != nil {
		return err
	}

	// 构建新的 crontab 内容
	newJob := fmt.Sprintf("%s %s", schedule, command)
	newCrontab := buildNewCrontab(currentCrontab, newJob)

	// 写入新的 crontab
	return writeCrontab(newCrontab)
}

// getCurrentCrontab 获取当前的crontab内容
func getCurrentCrontab() (string, error) {
	cmd := exec.Command("crontab", "-l")
	output, err := cmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok && exitError.ExitCode() == 1 {
			return "", nil // 没有现有的 crontab
		}
		return "", fmt.Errorf("获取当前 crontab 失败: %v", err)
	}
	return string(output), nil
}

// buildNewCrontab 构建新的crontab内容
func buildNewCrontab(current, newJob string) string {
	if current != "" && !strings.HasSuffix(current, "\n") {
		current += "\n"
	}
	return current + newJob + "\n"
}

// editCronJobs 编辑 cron 任务
func editCronJobs() error {
	if runtime.GOOS == "windows" {
		return editWindowsScheduledTasks()
	}

	cmd := exec.Command("crontab", "-e")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// monitorCronJobs 监控 cron 任务
func monitorCronJobs() error {
	if runtime.GOOS == "windows" {
		return monitorWindowsScheduledTasks()
	}
	return monitorUnixCronJobs()
}

// monitorUnixCronJobs 监控Unix系统的cron任务
func monitorUnixCronJobs() error {
	fmt.Println("开始监控 cron 任务...")
	fmt.Println("按 Ctrl+C 停止监控")
	fmt.Println(separatorLine)

	// 尝试找到可用的日志文件
	logPaths := []string{
		"/var/log/cron",
		"/var/log/cron.log",
		"/var/log/syslog",
	}

	activeLogPath := findActiveLogPath(logPaths)
	if activeLogPath == "" {
		fmt.Println("警告: 未找到 cron 日志文件，尝试监控进程...")
		return monitorCronProcesses()
	}

	return tailCronLog(activeLogPath)
}

// findActiveLogPath 查找可用的日志文件路径
func findActiveLogPath(paths []string) string {
	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	return ""
}

// tailCronLog 追踪cron日志文件
func tailCronLog(logPath string) error {
	cmd := exec.Command("tail", "-f", logPath)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("创建管道失败: %v", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("启动监控失败: %v", err)
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(strings.ToLower(line), "cron") {
			fmt.Printf("[%s] %s\n", time.Now().Format(cronTimeFormat), line)
		}
	}

	return cmd.Wait()
}

// deleteCronJob 删除指定行的 cron 任务
func deleteCronJob(lineNum int) error {
	if runtime.GOOS == "windows" {
		return deleteWindowsScheduledTask(lineNum)
	}
	return deleteUnixCronJob(lineNum)
}

// deleteUnixCronJob 删除Unix系统的cron任务
func deleteUnixCronJob(lineNum int) error {
	// 获取当前 crontab
	currentCrontab, err := getCurrentCrontab()
	if err != nil {
		return fmt.Errorf("获取 cron 任务失败: %v", err)
	}

	lines := strings.Split(currentCrontab, "\n")
	if lineNum < 1 || lineNum > len(lines) {
		return fmt.Errorf("无效的行号: %d (总共 %d 行)", lineNum, len(lines))
	}

	// 删除指定行
	newLines := append(lines[:lineNum-1], lines[lineNum:]...)
	newCrontab := strings.Join(newLines, "\n")

	return writeCrontab(newCrontab)
}

// =============================================================================
// 辅助函数
// =============================================================================

// writeCrontab 写入新的 crontab
func writeCrontab(content string) error {
	cmd := exec.Command("crontab", "-")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("创建输入管道失败: %v", err)
	}

	go func() {
		defer stdin.Close()
		stdin.Write([]byte(content))
	}()

	return cmd.Run()
}

// isValidCronExpression 简单验证 cron 表达式
func isValidCronExpression(expr string) bool {
	parts := strings.Fields(expr)
	return len(parts) == 5
}

// monitorCronProcesses 监控 cron 进程
func monitorCronProcesses() error {
	fmt.Println("监控 cron 相关进程...")

	for {
		cmd := exec.Command("ps", "aux")
		output, err := cmd.Output()
		if err != nil {
			return fmt.Errorf("获取进程列表失败: %v", err)
		}

		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if strings.Contains(strings.ToLower(line), "cron") {
				fmt.Printf("[%s] %s\n", time.Now().Format(cronTimeFormat), line)
			}
		}

		time.Sleep(5 * time.Second)
	}
}

// =============================================================================
// Windows 相关功能实现
// =============================================================================

// listWindowsScheduledTasks 列出 Windows 计划任务
func listWindowsScheduledTasks() error {
	fmt.Println("Windows 计划任务:")
	fmt.Println(separatorLine)

	cmd := exec.Command("schtasks", "/query", "/fo", "table")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("获取计划任务失败: %v", err)
	}

	fmt.Print(string(output))
	return nil
}

// addWindowsScheduledTask 添加 Windows 计划任务
func addWindowsScheduledTask(schedule, command string) error {
	// 将 cron 表达式转换为 Windows 计划任务格式
	// 这是一个简化版本，实际需要更复杂的转换逻辑
	taskName := fmt.Sprintf("CronTask_%d", time.Now().Unix())

	cmd := exec.Command("schtasks", "/create", "/tn", taskName, "/tr", command, "/sc", "daily")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("创建 Windows 计划任务失败: %v", err)
	}

	fmt.Printf("创建 Windows 计划任务成功: %s\n", taskName)
	return nil
}

// editWindowsScheduledTasks 编辑 Windows 计划任务
func editWindowsScheduledTasks() error {
	fmt.Println("请使用 Windows任务计划程序 (taskschd.msc) 编辑计划任务")

	cmd := exec.Command("cmd", "/c", "start", "taskschd.msc")
	return cmd.Run()
}

// monitorWindowsScheduledTasks 监控 Windows 计划任务
func monitorWindowsScheduledTasks() error {
	fmt.Println("监控 Windows 计划任务日志...")
	fmt.Println("请查看事件查看器中的任务调度程序日志")

	cmd := exec.Command("cmd", "/c", "start", "eventvwr.msc")
	return cmd.Run()
}

// deleteWindowsScheduledTask 删除 Windows 计划任务
func deleteWindowsScheduledTask(lineNum int) error {
	return fmt.Errorf("Windows 计划任务删除需要任务名称，请使用任务计划程序手动删除")
}
