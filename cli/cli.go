package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/orzation/bobibo"
)

var (
	version string

	gif       bool
	inverse   bool
	scale     float64
	threshold int
)

func init() {
	flag.BoolVar(&gif, "g", false, "enable gif mode.")
	flag.BoolVar(&inverse, "v", false, "inverse the colors.")
	flag.Float64Var(&scale, "s", 0.5, "scale the size of arts. range: (0, +).")
	flag.IntVar(&threshold, "t", -1, "set the threshold of binarization. range: [-1, 255], -1 means gen by OTSU.")
}

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "Usage: bobibo [OPTION]... PARTERNS [FILE]...")
		fmt.Fprintln(os.Stderr, "Try 'bobibo --help' for more information.")
		os.Exit(1)
	}

	opt := args[0]
	var imgFile *os.File

	switch opt {
	case "version":
		fmt.Printf("BoBiBo %s :P\n", version)
		return
	default:
		f, err := os.Open(opt)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Open image error: ", err.Error())
		}
		imgFile = f
	}
	defer imgFile.Close()
	arts, err := bobibo.BoBiBo(
		imgFile, gif, inverse,
		bobibo.ScaleOpt(scale),
		bobibo.ThresholdOpt(threshold))
	if err != nil {
		fmt.Fprintln(os.Stderr, "Bobibo error: ", err.Error())
		imgFile.Close()
		os.Exit(1)
	}

	err = printArts(arts, gif)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Print error: ", err.Error())
		imgFile.Close()
		os.Exit(1)
	}
}
