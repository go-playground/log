package log

import "sync"

// byteArrayPool represents a reusable byte pool. It is a centralized global instance for this package and can be
// accessed by calling log.BytePool(). It is intended to be used by Handlers.
type byteArrayPool struct {
	pool *sync.Pool
}

func (p *byteArrayPool) Get() *Buffer {
	return p.pool.Get().(*Buffer)
}

func (p *byteArrayPool) Put(buff *Buffer) {
	buff.B = buff.B[:0]
	p.pool.Put(buff)
}

// Buffer is a mere wrapper for a byte slice. It is intended to be used by Handlers.
type Buffer struct {
	B []byte
}
