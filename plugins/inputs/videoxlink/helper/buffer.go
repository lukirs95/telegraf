package helper

import (
	"errors"
	"sync"
)

const minCapacity = 16

var errEmptyBuffer = errors.New("buffer is empty")

type RingBuffer[T any] struct {
	bufMu sync.Mutex
	buf []T
	count uint

	
	rPtr uint
	wPtr uint
}

func NewRingBuffer[T any]() *RingBuffer[T] {
	return &RingBuffer[T]{
		buf: make([]T, minCapacity),
		count: 0,

		rPtr: 0,
		wPtr: 0,
	}
}

func (r *RingBuffer[T]) Cap() int {
	return len(r.buf)
}

func (r *RingBuffer[T]) Size() int {
	return int(r.count)
}

func (r *RingBuffer[T]) PushBack(elem T) {
	r.bufMu.Lock()
	defer r.bufMu.Unlock()
	r.growIfFull()

	r.buf[r.wPtr] = elem
	r.wPtr = r.next(r.wPtr)
	r.count++
}

func (r *RingBuffer[T]) PopFront() (T, error) {
	r.bufMu.Lock()
	defer r.bufMu.Unlock()
	if r.count == 0 {
		var zero T
		return zero, errEmptyBuffer
	}
	retVal := r.buf[r.rPtr]
	var zero T
	r.buf[r.rPtr] = zero
	r.rPtr = r.next(r.rPtr)
	r.count--

	return retVal, nil
}

func (r *RingBuffer[T]) next(i uint) uint {
	return (i + 1) & (uint(len(r.buf)) - 1) // bitwise modulus
}

func (r *RingBuffer[T]) growIfFull() {
	if r.count != uint(len(r.buf)) {
		return
	}
	r.resize(r.count << 1)
}

func (r *RingBuffer[T]) resize(newSize uint) {
	newBuf := make([]T, newSize)
	if r.wPtr > r.rPtr {
		copy(newBuf, r.buf[r.rPtr:r.wPtr])
	} else {
		n := copy(newBuf, r.buf[r.rPtr:])
		copy(newBuf[n:], r.buf[:r.wPtr])
	}

	r.rPtr = 0
	r.wPtr = r.count
	r.buf = newBuf
}