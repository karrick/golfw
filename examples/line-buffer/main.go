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
	copyBufSize = 32768
	lfwBufSize  = 16384
)

func main() {
	lf, err1 := golfw.NewWriteCloser(os.Stdout, lfwBufSize)
	if err1 != nil {
		bail(1, err1)
	}

	_, err1 = io.CopyBuffer(lf, os.Stdin, make([]byte, copyBufSize))
	err2 := lf.Close() // NOTE: Also closes underlying io.WriteCloser.

	if err1 != nil {
		bail(1, err1)
	}
	if err2 != nil {
		bail(1, err2)
	}
}

func bail(code int, err error) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", filepath.Base(os.Args[0]), err)
	os.Exit(code)
}
