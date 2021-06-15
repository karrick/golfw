package golfw

import (
	"io"
	"testing"

	"github.com/karrick/gorill"
)

func benchmarkIt(b *testing.B, iowc io.WriteCloser) {
	for i := 0; i < 1<<20; i++ {
		n, err := iowc.Write([]byte{byte(i & 0xFF)})
		if got, want := n, 1; got != want {
			b.Errorf("GOT: %v; WANT: %v", got, want)
		}
		ensureError(b, err)
	}
	err := iowc.Close()
	ensureError(b, err)
}

func BenchmarkDevNull(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchmarkIt(b, gorill.NopCloseWriter(io.Discard))
	}
}

func BenchmarkWriteCloser(b *testing.B) {
	for i := 0; i < b.N; i++ {
		lfwc, err := NewWriteCloser(gorill.NopCloseWriter(io.Discard), 128)
		ensureError(b, err)
		benchmarkIt(b, lfwc)
	}
}
