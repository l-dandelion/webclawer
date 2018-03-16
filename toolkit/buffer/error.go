package buffer

import (
	"errors"
)

var (
	ErrClosedBufferPool = errors.New("closed buffer pool")
	ErrClosedBuffer     = errors.New("closed buffer")
)
