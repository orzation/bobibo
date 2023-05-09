package img

import (
	"image"
)

type Pale = *image.Paletted

type Img struct {
	id    int
	value image.Image
}

func NewImg(id int, value image.Image) Img {
	return Img{id: id, value: value}
}

type Gray struct {
	id    int
	value [][]uint8
}

func (g Gray) Id() int {
	return g.id
}

func (g Gray) size() (int, int) {
	return len(g.value), len(g.value[0])
}
