package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "",
	Long:  "",
	Run:   runConfigCmd,
}

func init() {
	rootCmd.AddCommand(configCmd)
}

func runConfigCmd(_ *cobra.Command, _ []string) {
	viper.SetConfigName("env")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./conf")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("read config error: ", err)
	}

	log.Println(viper.AllSettings())
}
