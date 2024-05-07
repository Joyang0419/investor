package grpcx

import (
	"context"
	"fmt"
	"sync"

	"google.golang.org/grpc"

	"tools/errorx"
	"tools/logger"
)

type GrpcConnectionPool struct {
	conns            []*grpc.ClientConn
	maxConnectionNum int
	lock             sync.Mutex
}

func NewGrpcConnectionPool(
	ctx context.Context,
	serverAddr string,
	maxConnectionNum int,
	options ...grpc.DialOption,
) *GrpcConnectionPool {
	pool := &GrpcConnectionPool{
		conns:            make([]*grpc.ClientConn, 0, maxConnectionNum),
		maxConnectionNum: maxConnectionNum,
	}

	for i := 0; i < maxConnectionNum; i++ {
		conn, err := grpc.DialContext(
			ctx,
			serverAddr,
			options...,
		)
		if errorx.CheckErrorExist(err) {
			pool.CloseAllConnectionsOfPool() // CloseAllConnectionsOfPool all previously opened connections
			logger.Fatal("[GrpcConnectionPool] grpc.Dial err: %v", err)
		}
		pool.conns = append(pool.conns, conn)
	}

	return pool
}

// GetConnFromPool 從連線池中取得連線
func (p *GrpcConnectionPool) GetConnFromPool() (*grpc.ClientConn, error) {
	p.lock.Lock()
	defer p.lock.Unlock()

	if len(p.conns) > 0 {
		conn := p.conns[0]
		p.conns = p.conns[1:]
		return conn, nil
	}
	return nil, fmt.Errorf("[GrpcConnectionPool] GetConnFromPool err: no available connection")
}

// ReturnConnectionToPool 回收連線
func (p *GrpcConnectionPool) ReturnConnectionToPool(conn *grpc.ClientConn) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.conns = append(p.conns, conn)
}

// CloseAllConnectionsOfPool 關閉所有連線
func (p *GrpcConnectionPool) CloseAllConnectionsOfPool() {
	p.lock.Lock()
	defer p.lock.Unlock()
	for _, conn := range p.conns {
		if err := conn.Close(); err != nil {
			logger.Error("[GrpcConnectionPool] CloseAllConnectionsOfPool err: %v", err)
		}
	}
	p.conns = nil
}
