package viperx

import (
	"errors"

	"github.com/spf13/viper"

	"tools/reflectx"
)

// EnvSetIntoConfig set environment variables into config
// configName: config file name, example: "env"
// configType: config file type, example: "yaml"
// configPath: config file path, example: "./conf"
// config: config struct 必須是指標
func EnvSetIntoConfig(configName string, configType, configPath string, config any) error {
	if !reflectx.IsStructPtr(config) {
		return errors.New("config must be a struct pointer")
	}

	viper.SetConfigName(configName)
	viper.SetConfigType(configType)
	viper.AddConfigPath(configPath)

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return viper.Unmarshal(config)
}

func GetAllSettings() map[string]any {
	return viper.AllSettings()
}
