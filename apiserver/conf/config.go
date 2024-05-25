package conf

import (
	"log"

	"github.com/spf13/viper"
)

var Config = &struct {
	Server `yaml:"server"`
	Jwt    `yaml:"jwt"`
	Oauth2 `yaml:"oauth2"`
}{}

type Server struct {
	Port string `yaml:"port"`
}

type Jwt struct {
	Secret string `yaml:"secret"`
}

type Oauth2 struct {
	Google Google `yaml:"google"`
}

type Google struct {
	ClientId     string   `yaml:"clientID"`
	ClientSecret string   `yaml:"clientSecret"`
	RedirectUrl  string   `yaml:"redirectURL"`
	Scopes       []string `yaml:"scopes"`
}

func Init() {
	// viper get config from env.yaml
	viper.SetConfigName("env")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./conf")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("error reading config file, %s", err)
	}

	err := viper.Unmarshal(Config)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	// default value
	{
		viper.SetDefault("server.port", ":8080")

		viper.SetDefault("jwt.secret", []byte{})

		viper.SetDefault("oauth2.google.clientID", "")
		viper.SetDefault("oauth2.google.clientSecret", "")
		viper.SetDefault("oauth2.google.redirectURL", "")
		viper.SetDefault("oauth2.google.scopes", []string{})
	}
}
