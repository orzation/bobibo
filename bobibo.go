package bobibo

import (
	"errors"
	"image"
	"io"
	"runtime"

	"github.com/orzation/bobibo/img"
	u "github.com/orzation/bobibo/util"
)

type Params struct {
	Image     io.Reader
	Gif       bool
	Inverse   bool
	Scale     float64
	Threshold int
}

type Option func(p *Params) error

func ScaleOpt(scale float64) Option {
	return func(p *Params) error {
		if scale <= 0 {
			return errors.New("The Value of scale must be within (0, +).")
		}
		p.Scale = scale
		return nil
	}
}

func ThresholdOpt(thre int) Option {
	return func(p *Params) error {
		if thre < -1 || thre > 255 {
			return errors.New("The Value of threshold must be within [-1, 255].")
		}
		p.Threshold = thre
		return nil
	}
}

type Art struct {
	Content []string
	Delay   int
}

func BoBiBo(ima io.Reader, isGif, isInverse bool, opts ...Option) (<-chan Art, error) {
	params := &Params{
		Image:   ima,
		Gif:     isGif,
		Inverse: isInverse,
	}
	for _, opt := range opts {
		if err := opt(params); err != nil {
			return nil, err
		}
	}

	var maxCpu = runtime.NumCPU()

	// clone ResizeAndGray fn then sort it.
	ragFns := func(in <-chan img.Img) <-chan img.Gray {
		num := 1
		if cap(in) > 1 {
			// test
			num = maxCpu - 3
		}
		return u.SortChan(
			u.CloneChanFn(img.ResizeAndGray(params.Scale), num, in))
	}

	mix := u.Multiply(img.ArtotBin(params.Inverse),
		u.Multiply(img.BinotImg(params.Threshold),
			ragFns,
		))

	inStream, delays, err := analyzeImage(params)
	if err != nil {
		return nil, err
	}

	outStream := mix(inStream)
	wrap := wrapOut(delays)
	return wrap(outStream), nil
}

var wrapOut = func(delays []int) func(<-chan []string) <-chan Art {
	flag := true
	if delays == nil || len(delays) == 0 {
		flag = false
	}
	return u.GenChanFn(func(in <-chan []string, out chan<- Art) {
		cnt := 0
		for i := range in {
			if flag {
				out <- Art{Content: i, Delay: delays[cnt]}
				cnt++
			} else {
				out <- Art{Content: i, Delay: 0}
			}
		}
	})
}

func analyzeImage(params *Params) (<-chan img.Img, []int, error) {
	var delays []int
	var inChan <-chan img.Img
	if params.Gif {
		p, dls, err := img.LoadAGif(params.Image)
		if err != nil {
			return nil, nil, err
		}
		delays = dls
		inChan = newInStream(p...)
	} else {
		i, err := img.LoadAImage(params.Image)
		if err != nil {
			return nil, nil, err
		}
		inChan = newInStream(i)
	}
	return inChan, delays, nil
}

func newInStream[T image.Image](ims ...T) <-chan img.Img {
	in := make(chan img.Img, len(ims))
	defer close(in)
	for i, v := range ims {
		imgV := img.NewImg(i, v)
		in <- imgV
	}
	return in
}
