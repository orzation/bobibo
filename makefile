.PHONY: build block cpu mem

test:
	go test -bench=. -cpu=4 -blockprofile=block.pprof -cpuprofile=cpu.pprof -memprofile=mem.pprof
block: block.pprof
	go tool pprof -http=:9999 block.pprof
cpu: cpu.pprof
	go tool pprof -http=:9999 cpu.pprof
mem: mem.pprof
	go tool pprof -http=:9999 mem.pprof

