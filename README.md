# golfw

Go line feed writer

[![MIT license](https://img.shields.io/badge/License-MIT-blue.svg)](https://lbesson.mit-license.org/)
[![made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](http://golang.org)
[![GoDoc](https://godoc.org/github.com/karrick/golfw?status.svg)](https://godoc.org/github.com/karrick/golfw)
[![GoReportCard](https://goreportcard.com/badge/github.com/karrick/golfw)](https://goreportcard.com/report/github.com/karrick/golfw)

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
[lumberjack.Logger](https://github.com/natefinch/lumberjack) with
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
[hyperfine](https://github.com/sharkdp/hyperfine).

```
$ hyperfine --export-markdown ../../BENCHMARKS.md --prepare 'rm -f *.log' --warmup 10 \
            'cat 2600-h-1000.htm > stdout.log' \
            'cat 2600-h-1000.htm | ./log-rotator'
Benchmark #1: cat 2600-h-1000.htm > stdout.log
  Time (mean ± σ):      3.921 s ±  0.025 s    [User: 8.3 ms, System: 3840.0 ms]
  Range (min … max):    3.891 s …  3.964 s    10 runs

Benchmark #2: cat 2600-h-1000.htm | ./log-rotator
  Time (mean ± σ):      4.507 s ±  0.062 s    [User: 239.4 ms, System: 6198.0 ms]
  Range (min … max):    4.379 s …  4.595 s    10 runs

Summary
  'cat 2600-h-1000.htm > stdout.log' ran
    1.15 ± 0.02 times faster than 'cat 2600-h-1000.htm | ./log-rotator'
```

| Command | Mean [s] | Min [s] | Max [s] | Relative |
|:---|---:|---:|---:|---:|
| `cat 2600-h-1000.htm > stdout.log` | 3.911 ± 0.043 | 3.861 | 3.984 | 1.00 |
| `cat 2600-h-1000.htm \| ./log-rotator` | 4.510 ± 0.050 | 4.431 | 4.611 | 1.15 ± 0.02 |

Using mean run times, there is a 15% pipeline penalty using
golfw.WriteCloser on top of lumberjack.Logger compared to merely using
the `cat` UNIX utility, but shown below, the penalty from
golfw.WriteCloser is between 1% and 5%.

### golfw.WriteCloser compared to io.Discard

To determine how much overhead golfw.WriteCloser adds to a data
pipeline, I created a benchmark file `benchmark_test.go` to stream one
megabyte of sequential bytes directly to io.Discard, and another
benchmark that streams the same bytes through golfw.WriteCloser to
io.Discard. On my development system, streaming through
golfw.WriteCloser measured between 3% more efficient, and 4% less
efficient than merely streaming through io.Discard, depending on the
selected flush threshold.

```
$ go test -bench=.
goos: linux
goarch: amd64
pkg: github.com/karrick/golfw
cpu: AMD Ryzen Threadripper 3960X 24-Core Processor 
BenchmarkWriteCloser-48    	       4	 278117669 ns/op
BenchmarkDevNull-48        	       4	 286373981 ns/op
PASS
ok  	github.com/karrick/golfw	4.503s
```
