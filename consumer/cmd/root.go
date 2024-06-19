package cmd

import (
	"consumer/conf"

	"github.com/spf13/cobra"

	"tools/logger"
	"tools/viperx"
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic("root cmd execute error")
	}
}

func init() {
	if err := viperx.EnvSetIntoConfig("env", "yaml", "./conf", &conf.Config); err != nil {
		logger.Fatal("viperx.EnvSetIntoConfig err: %v", err)
	}
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(serverCmd)
}

var rootCmd = &cobra.Command{
	Short: "Root short description",
	Long:  "Root long description",
}
