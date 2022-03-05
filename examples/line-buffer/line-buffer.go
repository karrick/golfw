package main

// Read from standard input, buffers on newline while writing to
// standard output.
//
// This program is meant to serve as an example of how to use this
// library, and can be used as a benchmark to show the overhead of
// this program over the `cat` UNIX utility.

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/karrick/golfw"
)

const (
	copyBufSize = 512
	lfwBufSize  = 512
)

func main() {
	lfw, err := golfw.NewWriteCloser(os.Stdout, lfwBufSize)
	if err != nil {
		bail(1, err)
	}

	_, err = io.CopyBuffer(lfw, os.Stdin, make([]byte, copyBufSize))
	cerr := lfw.Close() // NOTE: Also closes underlying io.WriteCloser.

	if err != nil {
		bail(1, err)
	}
	if cerr != nil {
		bail(1, cerr)
	}
}

func bail(code int, err error) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", filepath.Base(os.Args[0]), err)
	os.Exit(code)
}
