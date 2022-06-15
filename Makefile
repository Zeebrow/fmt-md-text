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
Architecture: amd64
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

test: clean build 
	FMT_MD_TEXT_BINARY=build/fmt-md-text-$(GIT_HASH) go test -v .

package-deb: build
	mkdir -p dist/fmt-md-text/DEBIAN
	mkdir -p dist/fmt-md-text$(DEB_INSTALL_DIR)
	cp build/fmt-md-text dist/fmt-md-text$(DEB_INSTALL_DIR)/fmt-md-text
	touch dist/fmt-md-text/DEBIAN/control
	echo "$$DEBIAN_CONTROL" > dist/fmt-md-text/DEBIAN/control
	dpkg-deb --build dist/fmt-md-text
	cp dist/*.deb build/
      
clean:
	rm -rf dist/
	rm -rf build/*

.PHONY: build

