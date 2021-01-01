package utils

import "sync/atomic"

type AtomicInt int64

// 计数器加一
func (i *AtomicInt) Incr(n int64) {
	atomic.AddInt64((*int64)(i), n)
}

// 计数器加一并获得最新值
func (i *AtomicInt) IncrAndGet(n int64) int64 {
	return atomic.AddInt64((*int64)(i), n)
}

// 计数器当前值
func (i *AtomicInt) Get() int64 {
	return atomic.LoadInt64((*int64)(i))
}

// 计数器重置
func (i *AtomicInt) Reset() {
	atomic.StoreInt64((*int64)(i), 0)
}
