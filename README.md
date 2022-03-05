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

Here is a very simple example that streams from standard input to
standard output, but buffers input until at least 512 bytes read, and
breaks only on newline characters.

See a more practical example in the `examples/log-roller/` directory.

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

When running tests with benchmarks, I observe an approximate 8.6%
overhead when writting to WriteCloser which writes to io.Discard, over
just writing to io.Discard.

```Bash
go test -bench=.
goos: linux
goarch: amd64
pkg: github.com/karrick/golfw
BenchmarkDevNull-48        	       4	 273557676 ns/op
BenchmarkWriteCloser-48    	       4	 297179524 ns/op
PASS
ok  	github.com/karrick/golfw	4.593s
```

Using either of the example programs to run the benchmarks requires
the [hyperfine](https://github.com/sharkdp/hyperfine) program
somewhere on your PATH.

### War and Peace

I created a small line-buffer program in `examples/line-buffer/` that
can be used to benchmark the overhead that WriteCloser adds to a
pipeline.

cat is 1.20 times faster than gocat.
gocat is 1.04 times faster than line-buffer.

## Examples

### log-roller

I created a small log-roller program in `examples/log-roller/` that
wraps [lumberjack.Logger](https://github.com/natefinch/lumberjack)
with WriteCloser, and streams 1000 copies of War and Peace through the
UNIX `cat` utility and the resulting `log-roller` program.

```
$ cd examples/line-buffer
$ make bench
```

```
$ hyperfine --export-markdown ../../BENCHMARKS.md --prepare 'rm -f *.log' --warmup 10 \
            'cat 2600-h-1000.htm > stdout.log' \
            'cat 2600-h-1000.htm | ./log-roller'
Benchmark #1: cat 2600-h-1000.htm > stdout.log
  Time (mean ± σ):      3.921 s ±  0.025 s    [User: 8.3 ms, System: 3840.0 ms]
  Range (min … max):    3.891 s …  3.964 s    10 runs

Benchmark #2: cat 2600-h-1000.htm | ./log-roller
  Time (mean ± σ):      4.507 s ±  0.062 s    [User: 239.4 ms, System: 6198.0 ms]
  Range (min … max):    4.379 s …  4.595 s    10 runs

Summary
  'cat 2600-h-1000.htm > stdout.log' ran
    1.15 ± 0.02 times faster than 'cat 2600-h-1000.htm | ./log-roller'
```

| Command | Mean [s] | Min [s] | Max [s] | Relative |
|:---|---:|---:|---:|---:|
| `cat 2600-h-1000.htm > stdout.log` | 3.911 ± 0.043 | 3.861 | 3.984 | 1.00 |
| `cat 2600-h-1000.htm \| ./log-roller` | 4.510 ± 0.050 | 4.431 | 4.611 | 1.15 ± 0.02 |

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
