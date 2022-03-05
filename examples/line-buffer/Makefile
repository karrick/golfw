.PHONY: bench bench_buf_size bench_gocat bench_line_buffer clean

bench: bench_buf_size bench_gocat bench_line_buffer

bench_buf_size: gocat 2600-h-1000.htm
	hyperfine --style color --warmup 10 \
		-L buf_size 512,1024,2048,4096,8192,16384,32768,65536 \
		'./gocat {buf_size} < 2600-h-1000.htm > /dev/null'

bench_gocat: gocat 2600-h-1000.htm
	hyperfine --style color --warmup 10 \
		'./gocat 512 < 2600-h-1000.htm > /dev/null' \
		'cat < 2600-h-1000.htm > /dev/null'

bench_line_buffer: line-buffer gocat 2600-h-1000.htm
	hyperfine --style color --export-markdown BENCHMARKS.md --warmup 10 \
		'./line-buffer < 2600-h-1000.htm > /dev/null' \
		'./gocat 512 < 2600-h-1000.htm > /dev/null'

clean:
	rm -f gocat line-buffer 2600-h*.htm

2600-h.htm:
	curl -LOC - https://www.gutenberg.org/files/2600/2600-h/2600-h.htm

2600-h-1000.htm: 2600-h.htm
	for i in $$(seq 1000); do cat 2600-h.htm >> $@ ; done

gocat: gocat.go
	go build -o $@ gocat.go

line-buffer: line-buffer.go
	go build -o $@ line-buffer.go
