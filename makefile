VERSION=V1.4.0
EXE=bobibo
DESTDIR :=

.PHONY: default build install uninstall block cpu mem
default: build

build: cli/cli.go
	go build -ldflags="-X 'main.version=$(VERSION)' -s -w" -o $(EXE)
	@echo Build Success !!!

install: $(EXE)
	install -Dm755 $(EXE) $(DESTDIR)/usr/bin/$(EXE)
	@echo install Success !!!

uninstall:
	rm -f $(DESTDIR)/usr/bin/$(EXE)
	@echo uninstall Success !!!

test:
	go test -bench=. -cpu=4 -blockprofile=block.pprof -cpuprofile=cpu.pprof -memprofile=mem.pprof

block: block.pprof
	go tool pprof -http=:9999 block.pprof

cpu: cpu.pprof
	go tool pprof -http=:9999 cpu.pprof

mem: mem.pprof

	go tool pprof -http=:9999 mem.pprof

