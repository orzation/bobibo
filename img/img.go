package img

import (
	"image"
	"image/gif"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"math"
	"strings"

	u "github.com/orzation/bobibo/util"
	xdraw "golang.org/x/image/draw"
)

type Img = image.Image
type Pale = *image.Paletted

// use braille chars to draw arts.
var brailleMap = []rune("⠀⢀⡀⣀⠠⢠⡠⣠⠄⢄⡄⣄⠤⢤⡤⣤" +
	"⠐⢐⡐⣐⠰⢰⡰⣰⠔⢔⡔⣔⠴⢴⡴⣴⠂⢂⡂⣂⠢⢢⡢⣢⠆⢆⡆⣆⠦⢦⡦⣦⠒⢒⡒⣒⠲⢲⡲⣲⠖⢖⡖⣖⠶⢶⡶⣶" +
	"⠈⢈⡈⣈⠨⢨⡨⣨⠌⢌⡌⣌⠬⢬⡬⣬⠘⢘⡘⣘⠸⢸⡸⣸⠜⢜⡜⣜⠼⢼⡼⣼⠊⢊⡊⣊⠪⢪⡪⣪⠎⢎⡎⣎⠮⢮⡮⣮⠚⢚⡚⣚⠺⢺⡺⣺⠞⢞⡞⣞⠾⢾⡾⣾" +
	"⠁⢁⡁⣁⠡⢡⡡⣡⠅⢅⡅⣅⠥⢥⡥⣥⠑⢑⡑⣑⠱⢱⡱⣱⠕⢕⡕⣕⠵⢵⡵⣵⠃⢃⡃⣃⠣⢣⡣⣣⠇⢇⡇⣇⠧⢧⡧⣧⠓⢓⡓⣓⠳⢳⡳⣳⠗⢗⡗⣗⠷⢷⡷⣷" +
	"⠉⢉⡉⣉⠩⢩⡩⣩⠍⢍⡍⣍⠭⢭⡭⣭⠙⢙⡙⣙⠹⢹⡹⣹⠝⢝⡝⣝⠽⢽⡽⣽⠋⢋⡋⣋⠫⢫⡫⣫⠏⢏⡏⣏⠯⢯⡯⣯⠛⢛⡛⣛⠻⢻⡻⣻⠟⢟⡟⣟⠿⢿⡿⣿")

// loading an image, only support png and jpeg.
// if pass a gif, return the first embedded image.
// if there is any thing wrong, panic.
func LoadAImage(f io.Reader) (Img, error) {
	i, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	return i, nil
}

// loading a gif, return arrays of image.
func LoadAGif(f io.Reader) ([]Pale, []int, error) {
	g, err := gif.DecodeAll(f)
	if err != nil {
		return nil, nil, err
	}
	return g.Image, g.Delay, nil
}

// resizing the image with scale value, it won't change the ratio.
// return a stream chan function.
var ResizeAndGray = func(scale float64) func(<-chan Img) <-chan *image.Gray {
	return u.GenChanFunc(func(in <-chan Img, out chan<- *image.Gray) {
		for i := range in {
			dx := int(math.Floor(scale * float64(i.Bounds().Dx())))
			dy := int(math.Floor(scale * float64(i.Bounds().Dy())))
			dst := image.NewGray(image.Rect(0, 0, dx, dy))
			xdraw.NearestNeighbor.Scale(dst, dst.Rect, i, i.Bounds(), xdraw.Over, nil)
			out <- dst
		}
	})
}

// turning image to 2d binary matrix.
// use threshold to adjust the binarization.
var BinotImg = func(threshold int) func(<-chan *image.Gray) <-chan [][]bool {
	return u.GenChanFunc(func(in <-chan *image.Gray, out chan<- [][]bool) {
		for im := range in {
			out <- img2bin(im, &threshold)
		}
	})
}

func img2bin(im *image.Gray, th *int) [][]bool {
	if *th < 0 || *th > 255 {
		*th = int(otsu(im))
	}
	dx, dy := im.Bounds().Dx(), im.Bounds().Dy()
	reB := make([][]bool, dy)
	for i := range reB {
		reB[i] = make([]bool, dx)
		for j := range reB[i] {
			r, _, _, _ := im.At(j, i).RGBA()
			reB[i][j] = uint8(r>>8) >= uint8(*th)
		}
	}
	return reB
}

// return the best threshold to binarize.
func otsu(im *image.Gray) uint8 {
	var threshold int = 0
	const grayScale = 256
	var u float32
	var w0, u0 float32

	dx, dy := im.Bounds().Dx(), im.Bounds().Dy()
	hist := make([]float32, grayScale)
	sumPixel := dx * dy

	for i := 0; i < dy; i++ {
		for j := 0; j < dx; j++ {
			r, _, _, _ := im.At(j, i).RGBA()
			hist[uint8(r>>8)]++
		}
	}

	for i := range [grayScale]struct{}{} {
		hist[i] *= 1.0 / float32(sumPixel)
		u += float32(i) * hist[i]
	}
	var sigma float32
	for t := range [grayScale]struct{}{} {
		w0 += hist[t]
		u0 += float32(t) * hist[t]
		if w0 == 0 || 1-w0 == 0 {
			continue
		}

		tmp := u0 - u*w0
		tmp = tmp * tmp / (w0 * (1 - w0))
		if tmp >= sigma {
			sigma = tmp
			threshold = t
		}
	}
	return uint8(threshold)
}

// turning 2d binary matrix to string array.
// whether reverse color.
var ArtotBin = func(w bool) func(<-chan [][]bool) <-chan []string {
	return u.GenChanFunc(func(in <-chan [][]bool, out chan<- []string) {
		for e := range in {
			out <- bin2art(e, w)
		}
	})
}

func bin2art(bin [][]bool, isWhite bool) []string {
	dy, dx := len(bin)/4, len(bin[0])/2
	bufStr := make([]strings.Builder, dy)
	resStr := make([]string, dy)
	for i := 0; i < dy; i++ {
		for j := 0; j < dx; j++ {
			bufStr[i].WriteRune(cell(i, j, bin, isWhite))
		}
		resStr[i] = bufStr[i].String()
	}
	return resStr
}

func cell(y, x int, bin [][]bool, isWhite bool) rune {
	var reByte uint8 = 0
	for i := 0; i < 4; i++ {
		for j := 0; j < 2; j++ {
			if bin[y*4+i][x*2+j] {
				reByte += 1
			}
			if i != 3 || j != 1 {
				reByte <<= 1
			}
		}
	}
	if isWhite {
		reByte = ^reByte
	}
	return brailleMap[reByte]
}
