
# 讓protos產生pb.go
GenProtos:
	protoc --go_out=. ./protos/*/*.proto

# 啟動Gateway服務
RunGateway:
	cd gateway && go run main.go gateway
