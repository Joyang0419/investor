package cmd

import "github.com/spf13/cobra"

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic("root cmd execute error")
	}
}

func init() {
	// cmd setting...

}

var rootCmd = &cobra.Command{
	Short: "Root short description",
	Long:  "Root long description",
}
