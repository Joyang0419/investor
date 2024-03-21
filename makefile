# 讓protos產生pb.go
GenProtos:
	protoc --go_out=. --go-grpc_out=. ./protos/*/*.proto

# 啟動Gateway服務
RunGateway:
	cd gateway && go run main.go server

# 啟動micro_stock_price
RunMicroStockPrice:
	cd micro_stock_price && go run main.go server

# 啟動dev/build檔的dev docker compose yaml
UpDevInfra:
	cd build/dev && docker-compose up -d

# 關閉dev/build檔的dev docker compose yaml
DownDevInfra:
	cd build/dev && docker-compose down -v
