package util

// to make a stream chan function
func GenChanFunc[T any, E any](logic func(in <-chan T, out chan<- E)) func(<-chan T) <-chan E {
	return func(inChan <-chan T) <-chan E {
		outChan := make(chan E, len(inChan))
		go func() {
			defer close(outChan)
			logic(inChan, outChan)
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
