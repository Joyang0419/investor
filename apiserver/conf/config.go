package conf

var Config struct {
	Server struct {
		Port int
	}
	Jwt struct {
		Secret string
	}
	Oauth2 struct {
		Google struct {
			ClientId     string
			ClientSecret string
			RedirectUrl  string
			Scopes       []string
		}
	}
	GrpcServer struct {
		MicroAuth GrpcInfo
	}
}

type GrpcInfo struct {
	Domain           string
	Port             int
	MaxConnectionNum int
}
