package utils

import "sync/atomic"

type AtomicInt int64

func (i *AtomicInt) Incr(n int64) {
	atomic.AddInt64((*int64)(i), n)
}

func (i *AtomicInt) Get() int64 {
	return atomic.LoadInt64((*int64)(i))
}

func (i *AtomicInt) Reset() {
	atomic.StoreInt64((*int64)(i), 0)
}
