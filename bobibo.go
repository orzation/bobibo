package bobibo

import (
	"io"

	"github.com/orzation/bobibo/img"
	u "github.com/orzation/bobibo/util"
)

type Params struct {
	Image     io.Reader
	Gif       bool
	Reverse   bool
	Scale     float64
	Threshold int
}

type Option func(p *Params)

func BoBiBo(ima io.Reader, isGif, ifReverse bool, scale float64, threshold int, opts ...Option) (<-chan []string, error) {
	params := &Params{
		Image:     ima,
		Gif:       isGif,
		Reverse:   ifReverse,
		Scale:     scale,
		Threshold: threshold,
	}
	for _, opt := range opts {
		opt(params)
	}
	inStream := make(chan img.Img)

	mix := u.Multiply(img.ArtotBin(params.Reverse),
		u.Multiply(img.BinotImg(params.Threshold),
			u.Multiply(img.TurnGray,
				img.Resize(params.Scale))))

	outStream := mix(inStream)
	err := putStream(inStream, params)
	if err != nil {
		return nil, err
	}
	return outStream, nil
}

func putStream(in chan<- img.Img, params *Params) error {
	if params.Gif {
		p, err := img.LoadAGif(params.Image)
		if err != nil {
			return err
		}
		go func() {
			defer close(in)
			for _, v := range p {
				in <- v
			}
		}()
	} else {
		i, err := img.LoadAImage(params.Image)
		if err != nil {
			return err
		}
		go func() {
			defer close(in)
			in <- i
		}()
	}
	return nil
}
