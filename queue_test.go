package fastqueue

import (
	"sync"
	"testing"
)

var times = 20000
var paral = 100

func BenchmarkChannel(b *testing.B) {
	var wg sync.WaitGroup
	msg := make(chan interface{}, times*paral)
	wg.Add(paral)
	for i := 0; i < paral; i++ {
		go func(i int) {
			for j := 0; j < times; j++ {
				msg <- i
			}
		}(i)
	}
	for i := 0; i < paral; i++ {
		go func(i int) {
			for j := 0; j < times; j++ {
				<-msg
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func BenchmarkFastQueue(b *testing.B) {
	var wg sync.WaitGroup
	rb := NewQueue(uint64(times * paral))
	wg.Add(paral)
	for i := 0; i < paral; i++ {
		go func(i int) {
			for j := 0; j < times; j++ {
				rb.Push(i)
			}
		}(i)
	}
	for i := 0; i < paral; i++ {
		go func(i int) {
			for j := 0; j < times; j++ {
				rb.Pop()
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
}
