VERSION=V1.4.0
EXE=bobibo
DESTDIR :=

.PHONY: default build install uninstall
default: build

build: cli.go
	go build -ldflags="-X 'main.version=$(VERSION)' -s -w" -o $(EXE)
	@echo Build Success !!!

install: $(EXE)
	install -Dm755 $(EXE) $(DESTDIR)/usr/bin/$(EXE)
	@echo install Success !!!

uninstall:
	rm -f $(DESTDIR)/usr/bin/$(EXE)
	@echo uninstall Success !!!
