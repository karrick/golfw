package golfw

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func ensureError(tb testing.TB, err error, contains ...string) {
	tb.Helper()
	if len(contains) == 0 || (len(contains) == 1 && contains[0] == "") {
		if err != nil {
			tb.Fatalf("GOT: %v; WANT: %v", err, contains)
		}
	} else if err == nil {
		tb.Errorf("GOT: %v; WANT: %v", err, contains)
	} else {
		for _, stub := range contains {
			if stub != "" && !strings.Contains(err.Error(), stub) {
				tb.Errorf("GOT: %v; WANT: %q", err, stub)
			}
		}
	}
}

func ensureBuffer(tb testing.TB, got *bytes.Buffer, want string) {
	tb.Helper()
	if g, w := string(got.Bytes()), want; g != w {
		tb.Errorf("GOT: %q; WANT: %q", g, w)
	}
}

func ensureWrite(tb testing.TB, lfw *WriteCloser, p string) {
	tb.Helper()
	n, err := lfw.Write([]byte(p))
	if got, want := n, len(p); got != want {
		tb.Errorf("GOT: %v; WANT: %v", got, want)
	}
	ensureError(tb, err)
}

type wantState struct {
	buf                 string
	n                   int
	indexOfFinalNewline int
	isShortWrite        bool
}

func ensureWriteResponse(tb testing.TB, lfwc *WriteCloser, p string, state wantState) {
	tb.Helper()
	n, err := lfwc.Write([]byte(p))
	if got, want := n, state.n; got != want {
		tb.Errorf("BYTES WRITTEN: GOT: %v; WANT: %v", got, want)
	}
	if state.isShortWrite {
		ensureError(tb, err, io.ErrShortWrite.Error())
	} else {
		ensureError(tb, err)
	}
	if got, want := lfwc.buf, []byte(state.buf); !bytes.Equal(got, want) {
		tb.Errorf("GOT: %q; WANT: %q", got, want)
	}
	if got, want := lfwc.indexOfFinalNewline, state.indexOfFinalNewline; got != want {
		tb.Errorf("FINAL NEWLINE: GOT: %v; WANT: %v", got, want)
	}
}
