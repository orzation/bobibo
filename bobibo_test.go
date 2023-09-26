package bobibo

import (
	"fmt"
	"os"
	"testing"
)

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
					return
				}
			}
		}
	}
}
