package handler

import (
	"tools/grpcx"
)

type GrpcConnectionPools struct {
	MicroAuthGrpcConnPool *grpcx.GrpcConnectionPool
}

func NewGrpcConnectionPools(
	microAuthGrpcConnPool *grpcx.GrpcConnectionPool,
) GrpcConnectionPools {
	return GrpcConnectionPools{
		MicroAuthGrpcConnPool: microAuthGrpcConnPool,
	}
}
