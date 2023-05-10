package util

import (
	"sync"
)

// to make a stream channel function.
func GenChanFn[T any, E any](logic func(in <-chan T, out chan<- E)) func(<-chan T) <-chan E {
	return func(inChan <-chan T) <-chan E {
		outChan := make(chan E, cap(inChan))
		go func() {
			logic(inChan, outChan)
			close(outChan)
		}()
		return outChan
	}
}

// multiply two functions to one, right function first.
func Multiply[T, E, R any](f func(E) R, g func(T) E) func(T) R {
	return func(v T) R {
		return f(g(v))
	}
}

// cloning stream channel function means that more goroutines will work for it.
// and finally all stream will be faned in one channel.
func CloneChanFn[T, E any](fn func(<-chan T) <-chan E, num int, in <-chan T) <-chan E {
	out := make(chan E, cap(in))

	wg := sync.WaitGroup{}
	wg.Add(len(in))

	go func() {
		wg.Wait()
		close(out)
	}()

	for i := 0; i < num; i++ {
		go func() {
			for v := range fn(in) {
				out <- v
				wg.Done()
			}
		}()
	}

	return out
}

// use the id to sort.
type sorter interface {
	// start from zero will be eazier.
	Id() int
}

func SortChan[T sorter](in <-chan T) <-chan T {
	out := make(chan T, cap(in))
	go func() {
		defer close(out)
		order := make([]T, cap(in))
		for s := range in {
			order[s.Id()] = s
		}
		for _, v := range order {
			out <- v
		}
	}()
	return out
}
