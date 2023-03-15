VERSION=V1.1.0
EXE=bobibo
DESTDIR :=

.PHONY: default
default: build

.PHONY: build
build: cli.go
	go build -ldflags="-X 'main.version=$(VERSION)' -s -w" -o $(EXE)
	upx $(EXE)	# upx not found ? install it or remove this line.
	@echo Build Success !!!

.PHONY: install
install: $(EXE)
	install -Dm755 $(EXE) $(DESTDIR)/usr/bin/$(EXE)
	@echo install Success !!!

.PHONY: uninstall
uninstall:
	rm -f /usr/bin/$(EXE)
	@echo uninstall Success !!!
