PROG_NAME := drcat
GIT_HASH := $(shell git rev-parse --short HEAD)
GIT_HASH_LONG := $(shell git rev-parse HEAD)
BUILD_DATE := $(shell date -I)
GOARCH := amd64#amd64, 386, arm, ppc64
#GOOS := linux#linux, darwin, windows, netbsd

DEB_INSTALL_DIR := /usr/bin
define DEBIAN_CONTROL =
Package: $(PROG_NAME)
Version: 0.2
Provides: $(PROG_NAME)
Section: custom
Priority: optional
Architecture: $(GOARCH)
Essential: no
Installed-Size: 8192
Maintainer: zeebrow
Homepage: https://github.com/zeebrow/$(PROG_NAME)
Description: pretty-print markdown files in your console
endef
export DEBIAN_CONTROL

build:
	go install .
	go build -ldflags "\
		-X 'main.ProgramName=$(PROG_NAME)' \
		-X 'main.CommitHash=$(GIT_HASH_LONG)' \
		-X 'main.Version=$(GIT_HASH)' \
		-X 'main.BuildDate=$(BUILD_DATE)' \
		" \
		-o build/$(PROG_NAME) .
	DRCAT_BINARY_DIR=build go test -v

build-windows:
	go install .
	GOOS=windows go build -ldflags "\
		-X 'main.ProgramName=$(PROG_NAME)' \
		-X 'main.CommitHash=$(GIT_HASH_LONG)' \
		-X 'main.Version=$(GIT_HASH)' \
		-X 'main.BuildDate=$(BUILD_DATE)' \
		" \
		-o build/$(PROG_NAME).exe .
	DRCAT_BINARY_DIR=build go test -v

package-deb: build
	mkdir -p dist/$(PROG_NAME)/DEBIAN
	mkdir -p dist/$(PROG_NAME)$(DEB_INSTALL_DIR)
	cp build/$(PROG_NAME) dist/$(PROG_NAME)$(DEB_INSTALL_DIR)/$(PROG_NAME)
	touch dist/$(PROG_NAME)/DEBIAN/control
	echo "$$DEBIAN_CONTROL" > dist/$(PROG_NAME)/DEBIAN/control
	dpkg-deb --build dist/$(PROG_NAME)
	cp dist/*.deb build/

remove-deb:
	sudo apt -y remove $(PROG_NAME) 
reinstall-deb: clean remove-deb package-deb
	sudo apt install ./build/$(PROG_NAME).deb

clean:
	rm -rf dist/
	rm -rf build/*

.PHONY: build

