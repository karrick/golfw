package main

// Read from standard input, and writes to rotated logs.

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/karrick/golfw"
	"github.com/natefinch/lumberjack"
)

const (
	copyBufSize = 32768
	lfwBufSize  = 16384
)

func main() {
	lj := &lumberjack.Logger{
		Filename: filepath.Base(os.Args[0]) + ".log",
	}

	lf, err1 := golfw.NewWriteCloser(lj, lfwBufSize)
	if err1 != nil {
		bail(1, err1)
	}

	_, err1 = io.CopyBuffer(lf, os.Stdin, make([]byte, copyBufSize))
	err2 := lf.Close() // NOTE: Also closes underlying io.WriteCloser, namely lj.

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
