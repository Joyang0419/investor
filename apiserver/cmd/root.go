package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Short: "Root short description",
	Long:  "Root long description",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic("root cmd execute error")
	}
}

func init() {
	// cmd setting...

}
