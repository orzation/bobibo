package bobibo

import (
	"fmt"
	"os"
	"testing"
)

func TestBobibo(t *testing.T) {
	f, err := os.Open("./test.jpg")
	if err != nil {
		t.Error(err)
	}
	c, err2 := BoBiBo(f, false, false, 1.0, -1)
	if err2 != nil {
		panic(err2)
	}
	for e := range c {
		for _, v := range e {
			fmt.Println(v)
		}
	}

}
