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
		fmt.Println("Please input a path of image :P")
		fmt.Println("Use help to print options :P")
		fmt.Println("Use version to print version :P")
		os.Exit(1)
	}
	path := os.Args[1]
	if path == "help" {
		fmt.Println("Use -r to reverse the char color :P")
		fmt.Println("Use -g to enable gif analyzation, default: disable :P")
		fmt.Println("Use -s [d]{[0, +)} to set the scale of art :P")
		fmt.Println("Use -t [d]{[0, 255]} to set the threshold of binarization :P")
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
				fmt.Println(err.Error())
				fmt.Println("The range of scale must be [0, +)")
				os.Exit(1)
			}
			scale = f
		case "-t":
			i, err := strconv.ParseInt(os.Args[i+3], 10, 64)
			if err != nil {
				fmt.Println("The range of threshold must be [0, 255]")
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
