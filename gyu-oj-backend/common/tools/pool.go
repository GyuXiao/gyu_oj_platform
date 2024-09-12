package tools

import (
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"sync"
)

// simple goroutine pool

type SimpleGoPool struct {
	ch       chan struct{}
	wg       *sync.WaitGroup
	poolSize int
}

// 采用有缓冲channel实现,当channel满的时候阻塞

func NewSimpleGoPool(maxSize int) *SimpleGoPool {
	if maxSize <= 0 {
		panic("SimpleGoPool size must be positive")
	}
	return &SimpleGoPool{
		ch:       make(chan struct{}, maxSize),
		wg:       new(sync.WaitGroup),
		poolSize: maxSize,
	}
}

func (sgp *SimpleGoPool) Add(delta int) {
	if len(sgp.ch) == sgp.poolSize {
		logc.Errorv(context.Background(), "SimpleGoPool is full")
	}
	sgp.wg.Add(delta)
	for i := 0; i < delta; i++ {
		sgp.ch <- struct{}{}
	}
}

func (sgp *SimpleGoPool) Done() {
	<-sgp.ch
	sgp.wg.Done()
}

func (sgp *SimpleGoPool) Wait() {
	sgp.wg.Wait()
}
