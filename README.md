# golfw

Go line feed writer

## Description

Library provides WriteCloser structure that is an io.WriteCloser that
buffers output to ensure it only emits bytes to the underlying
io.WriteCloser on line feed boundaries.

## Example

See a full example in the examples/log-rotator/ directory.

```Go
func Example() error {
    // Flush completed lines to os.Stdout at least every 512 bytes.
    lf, err := golfw.NewWriteCloser(os.Stdout, 512)
    if err != nil {
        return err
    }

    // Give copy buffer some room.
    _, rerr := io.CopyBuffer(lf, os.Stdin, make([]byte, 4096))

    // Clean up
    cerr := lf.Close()
    if rerr == nil {
        return cerr
    }
    return rerr
}
```
