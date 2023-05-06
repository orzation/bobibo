package img

import (
	"image"
	"image/draw"
	"image/gif"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"math"

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
var Resize = func(scale float64) func(<-chan Img) <-chan Img {
	return u.GenChanFunc(func(in <-chan Img, out chan<- Img) {
		for i := range in {
			dx := int(math.Floor(scale * float64(i.Bounds().Dx())))
			dy := int(math.Floor(scale * float64(i.Bounds().Dy())))
			dst := image.NewRGBA(image.Rect(0, 0, dx, dy))
			xdraw.NearestNeighbor.Scale(dst, dst.Rect, i, i.Bounds(), xdraw.Over, nil)
			out <- dst
		}
	})
}

// turning gray.
var TurnGray = u.GenChanFunc(func(in <-chan Img, out chan<- Img) {
	for i := range in {
		dx, dy := i.Bounds().Dx(), i.Bounds().Dy()
		dst := image.NewGray(image.Rect(0, 0, dx, dy))
		draw.Draw(dst, dst.Bounds(), i, i.Bounds().Min, draw.Src)
		out <- dst
	}
})

// turning image to 2d binary matrix.
// use threshold to adjust the binarization.
var BinotImg = func(threshold int) func(<-chan Img) <-chan [][]bool {
	return u.GenChanFunc(func(in <-chan Img, out chan<- [][]bool) {
		for im := range in {
			out <- img2bin(im, threshold)
		}
	})
}

func img2bin(im Img, th int) [][]bool {
	if th < 0 || th > 255 {
		th = int(otsu(im))
	}
	dx, dy := im.Bounds().Dx(), im.Bounds().Dy()
	reB := make([][]bool, dy)
	for i := range reB {
		reB[i] = make([]bool, dx)
		for j := range reB[i] {
			r, _, _, _ := im.At(j, i).RGBA()
			reB[i][j] = uint8(r>>8) >= uint8(th)
		}
	}
	return reB
}

// return the best threshold to binarize.
func otsu(im Img) uint8 {
	var threshold uint8 = 0
	const grayScale = 256
	var u float32
	dx, dy := im.Bounds().Dx(), im.Bounds().Dy()
	grayPro := make([]float32, grayScale)
	pixelSum := dx * dy
	for i := 0; i < dy; i++ {
		for j := 0; j < dx; j++ {
			r, _, _, _ := im.At(j, i).RGBA()
			grayPro[uint8(r>>8)]++
		}
	}
	for i := 0; i < grayScale; i++ {
		grayPro[i] *= 1.0 / float32(pixelSum)
		u += float32(i) * grayPro[i]
	}
	var w1, u1, gmax float32
	for i := 0; i < grayScale; i++ {
		w1 += grayPro[i]
		u1 += float32(i) * grayPro[i]

		tmp := u1 - u*w1
		sigma := tmp * tmp / (w1 * (1 - w1))
		if sigma >= gmax {
			threshold = uint8(i)
			gmax = sigma
		}
	}
	return threshold
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
	reStr := make([]string, dy)
	for i := 0; i < dy; i++ {
		for j := 0; j < dx; j++ {
			reStr[i] += cell(i, j, bin, isWhite)
		}
	}
	return reStr
}

func cell(y, x int, bin [][]bool, isWhite bool) string {
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
	return string(brailleMap[reByte])
}
