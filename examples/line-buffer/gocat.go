package main

// gocat - read from standard input, write to standard output

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
)

func main() {
	if len(os.Args) != 2 {
		bail(2, errors.New("USAGE: gocat BUF_SIZE"))
	}

	bufSize, err := strconv.ParseUint(os.Args[1], 10, 64)
	if err != nil {
		bail(1, fmt.Errorf("cannot parse BUF_SIZE: %s", err))
	}

	_, err = io.CopyBuffer(os.Stdout, os.Stdin, make([]byte, int(bufSize)))
	if err != nil {
		bail(1, err)
	}
}

func bail(code int, err error) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", filepath.Base(os.Args[0]), err)
	os.Exit(code)
}
