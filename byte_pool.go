package log

import "sync"

// BytePool represents a reusable byte pool. It is a centralized global instance for this package and can be accessed by
// calling log.BytePool(). It is intended to be used by Handlers.
type ByteArrayPool struct {
	pool *sync.Pool
}

func (p *ByteArrayPool) Get() []byte {
	return p.pool.Get().([]byte)
}

func (p *ByteArrayPool) Put(b []byte) {
	p.pool.Put(b[:0])
}
