package golfw

// These structures copied from https://github.com/karrick/gorill
// project.

import "io"

// NopCloseWriter returns a structure that implements io.WriteCloser, but provides a no-op Close
// method.  It is useful when you have an io.Writer that you must pass to a method that requires an
// io.WriteCloser.  It is the counter-part to ioutil.NopCloser, but for io.Writer.
//
//   iowc := gorill.NopCloseWriter(iow)
//   iowc.Close() // does nothing
func NopCloseWriter(iow io.Writer) io.WriteCloser { return nopCloseWriter{iow} }

func (nopCloseWriter) Close() error { return nil }

type nopCloseWriter struct{ io.Writer }

// ShortWriter returns a structure that wraps an io.Writer, but returns io.ErrShortWrite when the
// number of bytes to write exceeds a preset limit.
func ShortWriter(w io.Writer, max int) io.Writer {
	return shortWriter{Writer: w, max: max}
}

func (s shortWriter) Write(data []byte) (int, error) {
	var short bool
	index := len(data)
	if index > s.max {
		index = s.max
		short = true
	}
	n, err := s.Writer.Write(data[:index])
	if short {
		return n, io.ErrShortWrite
	}
	return n, err
}

type shortWriter struct {
	io.Writer
	max int
}
