package bobibo

import (
	"errors"
	"io"

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

	inStream := make(chan img.Img)

	mix := u.Multiply(img.ArtotBin(params.Inverse),
		u.Multiply(img.BinotImg(params.Threshold),
			u.Multiply(img.TurnGray,
				img.Resize(params.Scale),
			)))

	outStream := mix(inStream)
	delays, err := putStream(inStream, params)
	wrap := wrapOut(delays)
	if err != nil {
		return nil, err
	}
	return wrap(outStream), nil
}

var wrapOut = func(delays []int) func(<-chan []string) <-chan Art {
	flag := true
	if delays == nil || len(delays) == 0 {
		flag = false
	}
	return u.GenChanFunc(func(out <-chan []string, wrapOut chan<- Art) {
		cnt := 0
		for o := range out {
			if flag {
				wrapOut <- Art{Content: o, Delay: delays[cnt]}
			} else {
				wrapOut <- Art{Content: o, Delay: 0}
			}
			cnt++
		}
	})
}

func putStream(in chan<- img.Img, params *Params) ([]int, error) {
	var delays []int
	if params.Gif {
		p, dls, err := img.LoadAGif(params.Image)
		if err != nil {
			return nil, err
		}
		delays = dls
		go inStream(in, p...)
	} else {
		i, err := img.LoadAImage(params.Image)
		if err != nil {
			return nil, err
		}
		go inStream(in, i)
	}
	return delays, nil
}

func inStream[T img.Img](in chan<- img.Img, ims ...T) {
	defer close(in)
	for _, v := range ims {
		in <- v
	}
}
