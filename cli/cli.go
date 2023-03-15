package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/orzation/bobibo"
)

var (
	version string
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("Please input a path of image.")
		fmt.Println("Or using help to print options.")
		fmt.Println("Or using version to print version.")
		fmt.Println(":P")
		os.Exit(1)
	}
	path := os.Args[1]
	if path == "help" {
		fmt.Println("Options:")
		fmt.Println("  -r                 reverse the char color.")
		fmt.Println("  -g                 enable gif analyzation, default: disable.")
		fmt.Println("  -s [d](0, +)       set the scale of art.  [default: 0.5]")
		fmt.Println("  -t [d][0, 255]     set the threshold of binarization.  [default: gen by ostu]")
		os.Exit(0)
	} else if path == "version" {
		fmt.Printf("BoBiBo %s :P\n", version)
		os.Exit(0)
	}
	f, err := os.Open(path)
	defer f.Close()

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var gif, rever bool
	var scale float64 = 0.5
	var threshold = -1
	for i, v := range os.Args[2:] {
		switch v {
		case "-g":
			gif = true
		case "-r":
			rever = true
		case "-s":
			f, err := strconv.ParseFloat(os.Args[i+3], 64)
			if err != nil {
				fmt.Println("The range of scale must at (0, +).")
				os.Exit(1)
			}
			if f == 0 {
				fmt.Println("The range of scale must at (0, +).")
				os.Exit(1)
			}
			scale = f
		case "-t":
			i, err := strconv.ParseInt(os.Args[i+3], 10, 64)
			if err != nil {
				fmt.Println("The range of threshold must at [0, 255].")
				os.Exit(1)
			}
			if i < 0 || i > 255 {
				fmt.Println("The range of threshold must at [0, 255].")
				os.Exit(1)
			}
			threshold = int(i)
		}
	}

	out, err := bobibo.BoBiBo(f, gif, rever, scale, threshold)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	for e := range out {
		for _, v := range e {
			fmt.Printf("\r%s\n", v)
		}
	}
}
