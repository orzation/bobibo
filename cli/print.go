package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"time"

	b "github.com/orzation/bobibo"
	"golang.org/x/term"
)

func printArts(arts <-chan b.Art, gifMode bool) error {
	if !gifMode {
		for art := range arts {
			for _, v := range art.Content {
				fmt.Println(v)
			}
		}
		return nil
	}

	artsBuffer := make([]b.Art, len(arts))
	for a := range arts {
		artsBuffer = append(artsBuffer, a)
	}

	fd := int(os.Stdout.Fd())
	errChan := make(chan error, 2)
	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, os.Interrupt)

	hideCursor()
	defer showCursor()

	go func() {
		for {
			tw, th, err := term.GetSize(fd)
			if err != nil {
				errChan <- err
				return
			}

			clearAll()
			for _, art := range artsBuffer {
				sw, sh := len([]rune(art.Content[0])), len(art.Content)
				posX, posY := (tw-sw)>>1, (th-sh)>>1
				if posX < 0 || posY < 0 || posX > tw || posY > th {
					errChan <- errors.New("Image size is too large, please zoom out and try again.")
					return
				}
				for offset, line := range art.Content {
					moveCursor(posY+offset, posX)
					clearLine()
					fmt.Printf("%s", line)
				}
				time.Sleep(time.Microsecond * time.Duration(art.Delay*10000))
			}
		}
	}()

	select {
	case <-interrupt:
		return nil
	case err := <-errChan:
		return err
	}
}

func moveCursor(y, x int) {
	fmt.Printf("\033[%d;%dH", y, x)
}

func clearAll() {
	fmt.Print("\033[2J")
}

func clearLine() {
	fmt.Print("\033[2K")
}

func hideCursor() {
	fmt.Print("\033[?25l")
}

func showCursor() {
	fmt.Print("\033[?25h")
}
