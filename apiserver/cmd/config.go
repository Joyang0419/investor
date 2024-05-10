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

// TODO To Eric delete this Command, because unnecessary
func runConfigCmd(_ *cobra.Command, _ []string) {
	viper.SetConfigName("env")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("read config error: ", err)
	}

	// TODO show all settings prettily
	log.Println(viper.AllSettings())
}
