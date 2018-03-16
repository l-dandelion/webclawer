package buffer

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
)

type Buffer interface {
	Cap() uint32
	Len() uint32
	Put(datum interface{}) (bool, error)
	Get() (interface{}, error)
	Close() bool
	Closed() bool
}

type myBuffer struct {
	ch          chan interface{}
	closed      uint32
	closingLock sync.RWMutex
}

func (buf *myBuffer) Cap() uint32 {
	return uint32(cap(buf.ch))
}

func (buf *myBuffer) Len() uint32 {
	return uint32(len(buf.ch))
}

func (buf *myBuffer) Put(datum interface{}) (ok bool, err error) {
	buf.closingLock.RLock()
	defer buf.closingLock.RUnlock()
	if buf.Closed() {
		return false, ErrClosedBuffer
	}
	select {
	case buf.ch <- datum:
		ok = true
	default:
		ok = false
	}
	return
}

func (buf *myBuffer) Get() (interface{}, error) {
	select {
	case datum, ok := <-buf.ch:
		if !ok {
			return nil, ErrClosedBuffer
		}
		return datum, nil
	default:
		return nil, nil
	}
}

func (buf *myBuffer) Close() bool {
	if atomic.CompareAndSwapUint32(&buf.closed, 0, 1) {
		buf.closingLock.Lock()
		defer buf.closingLock.Unlock()
		close(buf.ch)
		return true
	}
	return false
}

func (buf *myBuffer) Closed() bool {
	return atomic.LoadUint32(&buf.closed) == 1
}

func NewBuffer(size uint32) (Buffer, error) {
	if size == 0 {
		errMsg := fmt.Sprintf("illegal size for buffer: %d", size)
		return nil, errors.New(errMsg)
	}
	return &myBuffer{
		ch: make(chan interface{}, size),
	}, nil
}
