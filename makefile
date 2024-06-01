# 安裝GRPC套件
InstallGRPCPlugins:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

	export PATH="$${PATH}:$$(go env GOPATH)/bin"

# 讓protos產生pb.go
GenProtos:
	protoc --go_out=. --go-grpc_out=. ./protos/*.proto

# 產生GraphQL schema
GenGraphQL:
	cd apiserver/graphql && go run github.com/99designs/gqlgen generate

# 啟動ApiServer服務
RunApiServer:
	cd apiserver && cp conf/env.yaml.example conf/env.yaml && go mod tidy && go run main.go server

# 啟動micro_auth服務
RunMicroAuth:
	cd micro_auth && cp conf/env.yaml.example conf/env.yaml && go mod tidy && go run main.go server

# 啟動micro_stock_price
RunMicroStockPrice:
	cd micro_stock_price && go mod tidy && go run main.go server

# 啟動Scheduler服務
RunScheduler:
	cd scheduler && go mod tidy && go run main.go scheduler

# 啟動dev/build檔的dev docker compose yaml
UpDevInfra:
	cd build/dev && docker-compose up -d

# 關閉dev/build檔的dev docker compose yaml
DownDevInfra:
	cd build/dev && docker-compose down -v

MigrateUpMySQL:
	cd build/dev && docker-compose up -d flyway

# 啟動micro_stock_price dailyPrice 腳本, 記得要先建置MYSQL資料庫並執行Flyway, 確認table: daily_price存在
RunMicroStockPriceDailyPrice:
	cd micro_stock_price && go mod tidy && go run main.go dailyPrice

MigrateUpMongoDB:
	cd build/dev && docker-compose up -d migrate-mongo
