package bobibo

import (
	"fmt"
	"os"
	"testing"
)

func TestBobibo(t *testing.T) {
	f, err := os.Open("./test.gif")
	defer f.Close()
	if err != nil {
		t.Error(err)
	}
	c, err2 := BoBiBo(f, false, false, ScaleOpt(0.25), ThresholdOpt(-1))
	if err2 != nil {
		panic(err2)
	}
	for e := range c {
		for _, v := range e.Content {
			fmt.Println(v)
		}
	}
}

func BenchmarkBobibo(b *testing.B) {
	f, err := os.Open("./test.gif")
	if err != nil {
		b.Error(err)
	}
	defer f.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		f.Seek(0, 0)
		b.StartTimer()
		arts, err := BoBiBo(f, true, false, ScaleOpt(1), ThresholdOpt(-1))
		if err != nil {
			b.Error(err)
		}
		for {
			select {
			case _, ok := <-arts:
				if !ok {
					goto loopOut
				}
			}
		}
	loopOut:
	}
}
