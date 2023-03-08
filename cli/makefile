OS = linux
ARCH = amd64
EXE = bobibo_$(OS)_$(ARCH)
IS_STATIC = 0
all: build

build: cli.go
	CGO_ENABLED=$(IS_STATIC) GOOS=$(OS) GOARCH=$(ARCH) go build -ldflags="-s -w" -o $(EXE)
	@echo Build Success !!!

compress: $(EXE)
	upx $(EXE) -o $(EXE)_cprs
	@echo Compress Success !!!
