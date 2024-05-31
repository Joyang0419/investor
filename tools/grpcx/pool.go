package grpcx

import (
	"context"
	"fmt"
	"slices"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"

	"tools/errorx"
	"tools/logger"
	"tools/slicex"
)

type GrpcConnectionPool struct {
	conns            []*grpc.ClientConn
	maxConnectionNum int
	cursorIdx        int
	lock             sync.Mutex
}

func NewGrpcConnectionPool(
	ctx context.Context,
	serverAddr string,
	maxConnectionNum int,
	options ...grpc.DialOption,
) *GrpcConnectionPool {
	defaultMaxConnectionNum := 3
	if maxConnectionNum == 0 {
		maxConnectionNum = defaultMaxConnectionNum
	}

	pool := &GrpcConnectionPool{
		conns:            make([]*grpc.ClientConn, 0, maxConnectionNum),
		maxConnectionNum: maxConnectionNum,
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	options = append(options, grpc.WithBlock())
	for i := 0; i < maxConnectionNum; i++ {
		conn, err := grpc.DialContext(
			ctx,
			serverAddr,
			options...,
		)
		if errorx.IsErrorExist(err) {
			logger.Error("[NewGrpcConnectionPool]grpc.DialContext err: %v, serverAddr: %s", err, serverAddr)
			break
		}
		pool.conns = append(pool.conns, conn)
	}

	return pool
}

// GetConnFromPool 從連線池中取得連線
func (p *GrpcConnectionPool) GetConnFromPool() (*grpc.ClientConn, error) {
	p.lock.Lock()
	defer p.lock.Unlock()

	if slicex.IsEmpty(p.conns) {
		return nil, fmt.Errorf("[GrpcConnectionPool] GetConnFromPool err: no available connection")
	}

	for idx := p.getNextCursorIdx(); idx < len(p.conns); idx++ {
		/*
			Idle：此狀態表示 ClientConn 當前閒置，沒有任何活動。在此狀態下，ClientConn 可以被用來發起 RPC 請求。
			Connecting：此狀態表示 ClientConn 正在嘗試與服務端建立連接。在此狀態下，新的 RPC 請求可能會被阻塞或失敗，具體取決於你的 gRPC 配置。
			Ready：此狀態表示 ClientConn 已經成功與服務端建立連接，並且可以被用來發起新的 RPC 請求。
			TransientFailure：此狀態表示 ClientConn 目前無法與服務端建立連接，但是正在嘗試重新連接。在此狀態下，新的 RPC 請求可能會被阻塞或失敗。
			Shutdown：此狀態表示 ClientConn 已經被關閉，無法再被用來發起新的 RPC 請求。
		*/
		conn := p.conns[idx]
		validState := []connectivity.State{connectivity.Idle, connectivity.Ready}
		if slices.Contains(validState, conn.GetState()) {
			p.cursorIdx = idx
			return conn, nil
		}
		// 開go routine 重新連線, 客戶端無需等待這個操作
		go func() {
			conn.ResetConnectBackoff()
		}()
	}
	return nil, fmt.Errorf("[GrpcConnectionPool] GetConnFromPool err: no available connection")
}

// getNextCursorIdx 取得下一個游標位置, 當下一個游標位置超過連線池長度時, 重置游標位置
func (p *GrpcConnectionPool) getNextCursorIdx() int {
	p.cursorIdx++
	if slicex.IsIdxInSlice(p.conns, p.cursorIdx) {
		return p.cursorIdx
	}
	return 0
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
