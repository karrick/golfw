# golfw

Go line feed writer

## Description

Library provides WriteCloser structure that is an io.WriteCloser that
buffers output to ensure it only emits bytes to the underlying
io.WriteCloser on line feed boundaries.

## Example

```Go
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

	lf, err := golfw.NewWriteCloser(lj, lfwBufSize)
	if err != nil {
		bail(err, 1)
	}

	_, rerr := io.CopyBuffer(lf, os.Stdin, make([]byte, copyBufSize))
	cerr := lf.Close() // NOTE: Also closes underlying io.WriteCloser, namely lj.

	if rerr != nil {
		bail(rerr, 1)
	}
	if cerr != nil {
		bail(cerr, 1)
	}
}

func bail(err error, code int) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", filepath.Base(os.Args[0]), err)
	os.Exit(code)
}
```
