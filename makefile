
# 讓protos產生pb.go
GenProtos:
	protoc --go_out=. ./protos/*/*.proto
