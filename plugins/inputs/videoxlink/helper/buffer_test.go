package helper

import (
	"errors"
	"sync"
	"testing"
)

var runs int = 100000000

func TestRingBuffer(t *testing.T) {
	r := NewRingBuffer[int]()
	var wg sync.WaitGroup

	wg.Add(1)
	go func ()  {
		defer wg.Done()
		for i := 0; i < runs; i++ {
			r.PushBack(i)
		}
	}()
	wg.Add(1)
	go func () {
		defer wg.Done()
		for i := 0; i< runs; i++ {
			r.PushBack(i)
		}
	}()
	wg.Add(1)
	emptyBuffer := 0
	go func () {
		defer wg.Done()
		read := 0
		for {
			if read == runs*2-2 {
				return
			}
			_, err := r.PopFront()
			if errors.Is(err, errEmptyBuffer) {
				emptyBuffer++
				continue
			}
			read++
		}
	}()

	wg.Wait()
	t.Logf("buffercap: %d", r.Cap())
	t.Logf("empty reads: %d", emptyBuffer)
}