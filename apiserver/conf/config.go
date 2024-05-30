package conf

var Config struct {
	Server struct {
		Port string `yaml:"port"`
	} `yaml:"server"`
	Jwt struct {
		Secret string `yaml:"secret"`
	} `yaml:"jwt"`
	Oauth2 struct {
		Google struct {
			ClientId     string   `yaml:"clientID"`
			ClientSecret string   `yaml:"clientSecret"`
			RedirectUrl  string   `yaml:"redirectURL"`
			Scopes       []string `yaml:"scopes"`
		}
	} `yaml:"oauth2"`
}
