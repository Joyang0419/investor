package service

import (
	"github.com/panjf2000/ants/v2"

	"tools/errorx"
	"tools/logger"
)

type TaskPool struct {
	pool *ants.PoolWithFunc // ants 工作池
	size int                // 工作池大小
	task *Task              // 任務模板
}

func NewTaskPool(poolSize int, task *Task) *TaskPool {
	fn := func(i interface{}) {
		t := i.(*Task)
		t.Do()
	}

	pool, err := ants.NewPoolWithFunc(
		poolSize,
		fn,
	)
	if errorx.IsErrorExist(err) {
		logger.Fatal("Failed to create task pool: %v", err)
	}

	return &TaskPool{
		pool: pool,
		size: poolSize,
		task: task,
	}
}

// Start 啟動 TaskPool，提交任務到工作池
func (tp *TaskPool) Start() {
	for i := 0; i < tp.size; i++ {
		if err := tp.pool.Invoke(tp.task); errorx.IsErrorExist(err) {
			logger.Error("[TaskPool][Start]pool.Invoke err: %v", err)
		}
	}
}

// Release 釋放 TaskPool 資源
func (tp *TaskPool) Release() {
	tp.pool.Release()
}
