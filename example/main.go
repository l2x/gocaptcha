package main

import (
	"bufio"
	"fmt"
	"github.com/l2x/gocaptcha"
	"io"
	"os"
)

func main() {
	capt := gocaptcha.New()
	f, txt, err := capt.Create()
	if err != nil {
		fmt.Println(err)
	}

	fo, err := os.Create("out.png")
	if err != nil {
		panic(err)
	}
	// close fo on exit and check for its returned error
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()
	// make a write buffer
	w := bufio.NewWriter(fo)

	// make a buffer to keep chunks that are read
	buf := make([]byte, 1024)
	for {
		// read a chunk
		n, err := f.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if n == 0 {
			break
		}

		// write a chunk
		if _, err := w.Write(buf[:n]); err != nil {
			panic(err)
		}
	}

	if err = w.Flush(); err != nil {
		panic(err)
	}

	fmt.Println("Wrote out.png OK. text is ", txt)
}
