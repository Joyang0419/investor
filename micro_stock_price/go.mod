module micro_stock_price

go 1.21.5

replace (
	protos v0.0.0 => ../protos
	tools v0.0.0 => ../tools
)

require (
	github.com/spf13/cobra v1.8.0
	google.golang.org/grpc v1.62.1
	protos v0.0.0
	tools v0.0.0
)

require (
	github.com/bytedance/sonic v1.11.3 // indirect
	github.com/chenzhuoyu/base64x v0.0.0-20230717121745-296ad89f973d // indirect
	github.com/chenzhuoyu/iasm v0.9.0 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/klauspost/cpuid/v2 v2.0.9 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	golang.org/x/arch v0.0.0-20210923205945-b76863e36670 // indirect
	golang.org/x/net v0.22.0 // indirect
	golang.org/x/sys v0.18.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240123012728-ef4313101c80 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
)
