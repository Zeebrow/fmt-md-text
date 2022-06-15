GIT_HASH := $(shell git rev-parse --short HEAD)
GIT_HASH_LONG := $(shell git rev-parse HEAD)
GOARCH := amd64 #amd64, 386, arm, ppc64
GOOS := linux #linux, darwin, windows, netbsd

build-linux-amd64:
	go install .
	GOOS=linux GOARCH=amd64 \
       go build \
			 -ldflags "-X 'main.Version=$(GIT_HASH_LONG)'" \
			 -o build/fmt-md-text-$(GIT_HASH)-linux-amd64 .

build-windows-amd64:
	go install .
	GOOS=windows GOARCH=arm \
			 go build -o build/fmt-md-text-$(GIT_HASH)-windows-amd64 .

install-zeebrow:
	go install .
	go build \
		-ldflags "-X 'main.Version=$(GIT_HASH_LONG)'" \
		-o fmt-md-text .
	FMT_MD_TEXT_BINARY=fmt-md-text go test -v .
	cp fmt-md-text $(HOME)/.local/bin/scripts

build:
	go install .
	go build \
		-ldflags "-X 'main.Version=$(GIT_HASH_LONG)'" \
		-o build/fmt-md-text-$(GIT_HASH) .

test:
	FMT_MD_TEXT_BINARY=build/fmt-md-text-$(GIT_HASH) go test -v .
      
clean:
	rm -rf build/*
      
.PHONY: build clean

