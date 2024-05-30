package cmd

import (
	"github.com/spf13/cobra"

	"tools/logger"
	"tools/viperx"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "",
	Long:  "",
	Run:   runConfigCmd,
}

func runConfigCmd(_ *cobra.Command, _ []string) {
	logger.Info("configCmd viperx.GetAllSettings(): %v", viperx.GetAllSettings())
}
