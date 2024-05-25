package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"tools/cronx"
	"tools/logger"
)

var schedulerCmd = &cobra.Command{
	Use:   "scheduler",
	Short: "",
	Long:  "",
	Run:   runSchedulerCmd,
}

func init() {
	rootCmd.AddCommand(schedulerCmd)
}

// 每5秒执行一次: */5 * * * * *

func runSchedulerCmd(_ *cobra.Command, _ []string) {
	logger.Info("scheduler start")
	cronManager := cronx.NewCronManager()

	cronManager.AddFunc("*/5 * * * * *", func() {})

	cronManager.Start() // 启动调度器

	sig := make(chan os.Signal, 1)                      // 创建一个信号chan
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM) // 监听SIGINT和SIGTERM信号

	fmt.Println("Cron scheduler started. Press Ctrl+C to exit.")
	<-sig // 阻塞，直到收到信号

	cronManager.Stop() // 停止调度器
	logger.Info("scheduler stopped")
}
