- makefile
```
# 安裝GRPC套件
make InstallGRPCPlugins

# 產生protos.pb.go
make GenProtos

# 啟動Gateway服務
make RunGateway

# 讓protos產生pb.go
make GenProtos

# 啟動dev/build檔的dev docker compose yaml
make UpDevInfra

# 關閉dev/build檔的dev docker compose yaml
make DownDevInfra
```