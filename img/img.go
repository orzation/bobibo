package img

import (
	"image"
	"image/gif"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"strings"

	u "github.com/orzation/bobibo/util"
)

// use braille chars to draw arts.
const braille = "⠀⢀⡀⣀⠠⢠⡠⣠⠄⢄⡄⣄⠤⢤⡤⣤" +
	"⠐⢐⡐⣐⠰⢰⡰⣰⠔⢔⡔⣔⠴⢴⡴⣴⠂⢂⡂⣂⠢⢢⡢⣢⠆⢆⡆⣆⠦⢦⡦⣦⠒⢒⡒⣒⠲⢲⡲⣲⠖⢖⡖⣖⠶⢶⡶⣶" +
	"⠈⢈⡈⣈⠨⢨⡨⣨⠌⢌⡌⣌⠬⢬⡬⣬⠘⢘⡘⣘⠸⢸⡸⣸⠜⢜⡜⣜⠼⢼⡼⣼⠊⢊⡊⣊⠪⢪⡪⣪⠎⢎⡎⣎⠮⢮⡮⣮⠚⢚⡚⣚⠺⢺⡺⣺⠞⢞⡞⣞⠾⢾⡾⣾" +
	"⠁⢁⡁⣁⠡⢡⡡⣡⠅⢅⡅⣅⠥⢥⡥⣥⠑⢑⡑⣑⠱⢱⡱⣱⠕⢕⡕⣕⠵⢵⡵⣵⠃⢃⡃⣃⠣⢣⡣⣣⠇⢇⡇⣇⠧⢧⡧⣧⠓⢓⡓⣓⠳⢳⡳⣳⠗⢗⡗⣗⠷⢷⡷⣷" +
	"⠉⢉⡉⣉⠩⢩⡩⣩⠍⢍⡍⣍⠭⢭⡭⣭⠙⢙⡙⣙⠹⢹⡹⣹⠝⢝⡝⣝⠽⢽⡽⣽⠋⢋⡋⣋⠫⢫⡫⣫⠏⢏⡏⣏⠯⢯⡯⣯⠛⢛⡛⣛⠻⢻⡻⣻⠟⢟⡟⣟⠿⢿⡿⣿"

var brailleMap = []rune(braille)

// loading an image, only support png and jpeg.
// if pass a gif, return the first embedded image.
// if there is any thing wrong, panic.
func LoadAImage(f io.Reader) (image.Image, error) {
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

// resize the image with scale value, using nearestNeighbor.
// return a stream chan function.
var ResizeAndGray = func(scale float64) func(<-chan Img) <-chan Gray {
	return u.GenChanFn(func(in <-chan Img, out chan<- Gray) {
		for i := range in {
			out <- grayNearestNeighbor(scale, i)
		}
	})
}

func grayNearestNeighbor(scale float64, src Img) Gray {
	DY, DX := src.value.Bounds().Dy(), src.value.Bounds().Dx()
	dy := int(scale * float64(DY))
	dx := int(scale * float64(DX))
	tgt := make([][]uint8, dy)

	for i := 0; i < dy; i++ {
		tgt[i] = make([]uint8, dx)
		for j := 0; j < dx; j++ {
			x, y := int((float64(j) / scale)), int((float64(i) / scale))
			r, g, b, _ := src.value.At(x, y).RGBA()
			grayColor := (299*r + 587*g + 114*b) / 1000
			tgt[i][j] = uint8(grayColor >> 8)
		}
	}
	return Gray{id: src.id, value: tgt}
}

// turning image to 2d binary matrix.
// use threshold to adjust the binarization.
var BinotImg = func(threshold int) func(<-chan Gray) <-chan [][]bool {
	return u.GenChanFn(func(in <-chan Gray, out chan<- [][]bool) {
		for im := range in {
			out <- img2bin(im, &threshold)
		}
	})
}

func img2bin(im Gray, th *int) [][]bool {
	if *th < 0 || *th > 255 {
		*th = int(otsu(im))
	}
	dy, dx := im.size()
	reB := make([][]bool, dy)
	for i := range reB {
		reB[i] = make([]bool, dx)
		for j := range reB[i] {
			grayValue := im.value[i][j]
			reB[i][j] = grayValue >= uint8(*th)
		}
	}
	return reB
}

// return the best threshold to binarize.
func otsu(im Gray) uint8 {
	var threshold int = 0
	const grayScale = 256
	var u float32
	var w0, u0 float32

	dy, dx := im.size()
	hist := make([]float32, grayScale)
	sumPixel := dy * dx

	for i := 0; i < dy; i++ {
		for j := 0; j < dx; j++ {
			grayValue := im.value[i][j]
			hist[grayValue]++
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
	return u.GenChanFn(func(in <-chan [][]bool, out chan<- []string) {
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
