package img

import (
	"fmt"
	"os"
	"testing"

	u "github.com/orzation/bobibo/util"
)

func TestXxx(t *testing.T) {
	in := make(chan Img)
	mix := u.Multiply(ArtotBin(false),
		u.Multiply(BinotImg(128),
			u.Multiply(TurnGray, Resize(0.50))))
	out := mix(in)
	f, _ := os.Open("../w.gif")
	i, err := LoadAGif(f)
	if err != nil {
		t.Error(err.Error())
	}
	f.Close()
	go func() {
		defer close(in)
		for _, p := range i {
			in <- p
		}
		// in <- i
	}()
	for e := range out {
		for _, v := range e {
			fmt.Println(v)
		}
	}
}
