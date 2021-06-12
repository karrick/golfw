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

## Benchmarks

### War and Peace

I created a small log-rotator program in
`examples/log-rotator/main.go` that uses composition to wrap
(lumberjack.Logger)[github.com/natefinch/lumberjack] with
golfw.WriteCloser, and streamed 1000 copies of War and Peace through
the UNIX `cat` utility and the resulting `log-rotator` program

I created an input file by creating a large 3.2 GiB file with 1000
copies of War and Peace:

```
$ cd examples/log-rotator
$ go build
$ wget https://www.gutenberg.org/files/2600/2600-h/2600-h.htm
$ for i in $(seq 1000); do cat 2600-h.htm >> 2600-h-1000.htm ; done
```

I ran the benchmarks using
(hyperfine)[https://github.com/sharkdp/hyperfine].

```
$ hyperfine --prepare 'rm -f *.log' --warmup 10 \
            'cat 2600-h-1000.htm > stdout.log' \
            'cat 2600-h-1000.htm | ./log-rotator'
Benchmark #1: cat 2600-h-1000.htm > stdout.log
  Time (mean ± σ):      3.949 s ±  0.042 s    [User: 6.8 ms, System: 3865.6 ms]
  Range (min … max):    3.898 s …  4.037 s    10 runs

Benchmark #2: cat 2600-h-1000.htm | ./log-rotator
  Time (mean ± σ):      4.482 s ±  0.050 s    [User: 237.0 ms, System: 6147.5 ms]
  Range (min … max):    4.373 s …  4.559 s    10 runs

Summary
  'cat 2600-h-1000.htm > stdout.log' ran
    1.13 ± 0.02 times faster than 'cat 2600-h-1000.htm | ./log-rotator'
```

Using mean run times, 4.482 seconds divided by 3.949 seconds is less
than 14% pipeline penalty using golfw.WriteCloser on top of
lumberjack.Logger compared to merely using the `cat` UNIX utility.

### golfw.WriteCloser compared to io.Discard

To determine how much overhead golfw.WriteCloser adds to a data
pipeline, I created a benchmark file in `benchmarks/benchmark_test.go`
to stream a megabyte of sequential bytes directly to io.Discard, and
another benchmark that streams the same bytes through
golfw.WriteCloser to the same io.Discard. On my development system,
streaming through golfw.WriteCloser adds between 1% and 4% overhead to
the pipeline.

```
$ go test -bench=.
goos: linux
goarch: amd64
pkg: github.com/karrick/golfw
cpu: AMD Ryzen Threadripper 3960X 24-Core Processor 
BenchmarkDevNull-48        	       4	 271099505 ns/op
BenchmarkWriteCloser-48    	       4	 278314647 ns/op
PASS
ok  	github.com/karrick/golfw	4.495s
```
