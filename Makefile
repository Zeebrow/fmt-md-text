GIT_HASH := $(shell git rev-parse --short HEAD)
GIT_HASH_LONG := $(shell git rev-parse HEAD)
BUILD_DATE := $(date -I)
GOARCH := amd64#amd64, 386, arm, ppc64
GOOS := linux#linux, darwin, windows, netbsd

DEB_INSTALL_DIR := /usr/bin
define DEBIAN_CONTROL =
Package: fmt-md-text
Version: 0.1
Provides: fmt-md-text
Section: custom
Priority: optional
Architecture: $(GOARCH)
Essential: no
Installed-Size: 8192
Maintainer: zeebrow
Homepage: https://github.com/zeebrow/fmt-md-text
Description: pretty-print markdown files in your console
endef
export DEBIAN_CONTROL

install-zeebrow:
	go install .
	go build \
		-ldflags "\
		-X 'main.Version=$(GIT_HASH_LONG)'\
		-X 'main.BuildDate=$(BUILD_DATE)'
		" \
		-o fmt-md-text .
	FMT_MD_TEXT_BINARY=fmt-md-text go test -v .
	cp fmt-md-text $(HOME)/.local/bin/scripts

build:
	go install .
	go build -ldflags "\
		-X 'main.Version=$(GIT_HASH_LONG)' \
		-X 'main.BuildDate=$(BUILD_DATE)' \
	" \
		-o build/fmt-md-text .

build-dev:
	go install .
	go build -ldflags "\
		-X 'main.Version=$(GIT_HASH_LONG)' \
		-X 'main.BuildDate=$(BUILD_DATE)' \
	" \
		-o build/fmt-md-text-$(GIT_HASH) .

test-dev: clean build-dev
	FMT_MD_TEXT_BINARY=build/fmt-md-text-$(GIT_HASH) go test -v .

test: clean build 
	FMT_MD_TEXT_BINARY=build/fmt-md-text go test -v .


package-tar:
	tar -czf fmt-md-text-$(GIT_HASH).tar.gz *.go go.mod go.sum README.md Makefile

package-deb: test
	mkdir -p dist/fmt-md-text/DEBIAN
	mkdir -p dist/fmt-md-text$(DEB_INSTALL_DIR)
	cp build/fmt-md-text dist/fmt-md-text$(DEB_INSTALL_DIR)/fmt-md-text
	touch dist/fmt-md-text/DEBIAN/control
	echo "$$DEBIAN_CONTROL" > dist/fmt-md-text/DEBIAN/control
	dpkg-deb --build dist/fmt-md-text
	cp dist/*.deb build/

build-release:
	go install .
	go build -ldflags "\
		-X 'main.Version=$(GIT_HASH_LONG)' \
		-X 'main.BuildDate=$(BUILD_DATE)' \
	" \
		-o build/fmt-md-text-$(GIT_HASH) .

release:
ifeq ($(strip $(VERSION)),)
	@echo set VERSION and try again
else
	@echo Relase version '$(VERSION)'
	@echo $(shell echo $(VERSION) | grep 'v\d\.d\.\d')
endif

uhhh:
ifeq ($(VERSION),undefined)
	VERSION = $(GIT_HASH)
endif
	@echo $(VERSION)
      
clean:
	rm -rf dist/
	rm -rf build/*

.PHONY: build

