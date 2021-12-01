line-buffer: main.go
	go build -o $@ $^

clean:
	rm -f line-buffer 2600-h*.htm

bench: line-buffer 2600-h-1000.htm
	hyperfine --export-markdown BENCHMARKS.md --warmup 10 \
		'cat < 2600-h-1000.htm > /dev/null' \
		'./line-buffer < 2600-h-1000.htm > /dev/null'

2600-h.htm:
	curl -LOC - https://www.gutenberg.org/files/2600/2600-h/2600-h.htm

2600-h-1000.htm: 2600-h.htm
	for i in $$(seq 1000); do cat $< >> $@ ; done

.PHONY: bench clean
