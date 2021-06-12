package golfw

import (
	"bytes"
	"fmt"
	"io"
)

// WriteCloser is an io.WriteCloser that buffers output to ensure it only emits
// bytes to the underlying io.WriteCloser on line feed boundaries.
type WriteCloser struct {
	buf                 []byte
	iowc                io.WriteCloser
	flushThreshold      int // flush on LF after buffer this size or larger
	indexOfFinalNewline int // -1 when no newlines in buf
}

// NewWriteCloser returns new WriteCloser with the specified flush
// threshold. Whenever the buffer is greater than the specified threshold, it
// flushes the buffer, up to and including the final LF byte, to the underlying
// io.WriteCloser.
//
//     func Example() error {
//         // Flush completed lines to os.Stdout at least every 512 bytes.
//         lf, err := golfw.NewWriteCloser(os.Stdout, 512)
//         if err != nil {
//             return err
//         }
//
//         // Give copy buffer some room.
//         _, rerr := io.CopyBuffer(lf, os.Stdin, make([]byte, 4096))
//
//         // Clean up
//         cerr := lf.Close()
//         if rerr == nil {
//             return cerr
//         }
//         return rerr
//     }
func NewWriteCloser(iowc io.WriteCloser, flushThreshold int) (*WriteCloser, error) {
	if flushThreshold <= 0 {
		return nil, fmt.Errorf("cannot create WriteCloser when flushThreshold less than or equal to 0: %d", flushThreshold)
	}
	return &WriteCloser{
		iowc:                iowc,
		flushThreshold:      flushThreshold,
		indexOfFinalNewline: -1,
	}, nil
}

// Close writes all data in its buffer to the underlying io.WriteCloser,
// including bytes without a trailing LF, then closes the underlying
// io.WriteCloser. This will either return any error caused by writing the bytes
// to the underlying io.WriteCloser, or an error caused by closing it. Use this
// method when done with a WriteCloser to prevent data loss.
func (lbf *WriteCloser) Close() error {
	_, we := lbf.iowc.Write(lbf.buf)
	lbf.buf = nil
	lbf.indexOfFinalNewline = -1
	ce := lbf.iowc.Close()
	lbf.iowc = nil
	if we == nil {
		return ce
	}
	return we
}

// flush flushes buffer to underlying io.WriteCloser, up to and including
// specified index.
func (lbf *WriteCloser) flush(olen, dlen, index int) (int, error) {
	nw, err := lbf.iowc.Write(lbf.buf[:index])
	if nw > 0 {
		nc := copy(lbf.buf, lbf.buf[nw:])
		lbf.buf = lbf.buf[:nc]
	}
	if err == nil {
		lbf.indexOfFinalNewline -= nw
		return dlen, nil
	}
	// nb is the number new bytes from p that got written to file.
	nb := nw - olen
	if nb < 0 {
		lbf.buf = lbf.buf[:-nb]
		nb = 0
	} else {
		lbf.buf = lbf.buf[:0]
	}
	lbf.indexOfFinalNewline = bytes.LastIndexByte(lbf.buf, '\n')
	return nb, err
}

// Write appends bytes from p to internal buffer, flushing buffer up to and
// including the final LF when buffer length exceeds programmed threshold.
func (lbf *WriteCloser) Write(p []byte) (int, error) {
	olen := len(lbf.buf)
	lbf.buf = append(lbf.buf, p...)

	if finalIndex := bytes.LastIndexByte(p, '\n'); finalIndex >= 0 {
		lbf.indexOfFinalNewline = olen + finalIndex
	}

	if len(lbf.buf) <= lbf.flushThreshold || lbf.indexOfFinalNewline < 0 {
		// Either do not need to flush, or no newline in buffer
		return len(p), nil
	}

	// Buffer larger than threshold, and has LF: write everything up to and
	// including that LF.
	return lbf.flush(olen, len(p), lbf.indexOfFinalNewline+1)
}
