OS = linux
ARCH = amd64
EXE = bobibo_$(OS)_$(ARCH)
IS_STATIC = 0
VERSION=V1.0.1
all: build

build: cli.go
	CGO_ENABLED=$(IS_STATIC) GOOS=$(OS) GOARCH=$(ARCH) go build -ldflags="-X 'main.version=$(VERSION)' -s -w" -o $(EXE)
	@echo Build Success !!!

compress: $(EXE)
	upx $(EXE) -o $(EXE)_cprs
	@echo Compress Success !!!
