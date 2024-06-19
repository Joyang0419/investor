package cmd

import (
	"github.com/spf13/cobra"

	"scheduler/conf"
	"tools/logger"
	"tools/viperx"
)

var rootCmd = &cobra.Command{
	Short: "Root short description",
	Long:  "Root long description",
}

func init() {
	if err := viperx.EnvSetIntoConfig("env", "yaml", "./conf", &conf.Config); err != nil {
		logger.Fatal("viperx.EnvSetIntoConfig err: %v", err)
	}
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(serverCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic("root cmd execute error")
	}
}
