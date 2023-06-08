PROG_NAME := drcat
GIT_HASH := $(shell git rev-parse --short HEAD)
GIT_HASH_LONG := $(shell git rev-parse HEAD)
BUILD_DATE := $(shell date -I)
GOARCH := amd64#amd64, 386, arm, ppc64
#GOOS := linux#linux, darwin, windows, netbsd
OS := linux
VERSION := 0.0.1-dev
DEB_INSTALL_DIR := /usr/bin

define DEBIAN_CONTROL =
Package: $(PROG_NAME)
Version: $(VERSION)
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

define BASED_ADVERTISEMENT =
#!/bin/sh
set -e
echo "\n\n***************************************************\nGain access to member-exclusive offer, birthdays treat, and more perk like 20% off your next visit when you join Denny's Rewards program today please!\nvisit https://www.dennys.com/ to learn even more offerings as well!"

echo "\n***************************************************\n\n"
read -p 'Sign up now for this exclusive offer? Y/n ' accept_offer
return 0
endef
export BASED_ADVERTISEMENT

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

package-deb: build
	mkdir -p dist/$(PROG_NAME)/DEBIAN
	mkdir -p dist/$(PROG_NAME)$(DEB_INSTALL_DIR)
	cp build/$(PROG_NAME) dist/$(PROG_NAME)$(DEB_INSTALL_DIR)/$(PROG_NAME)
	touch dist/$(PROG_NAME)/DEBIAN/control
	touch dist/$(PROG_NAME)/DEBIAN/preinst
	chmod 775 dist/$(PROG_NAME)/DEBIAN/preinst
	echo "$$DEBIAN_CONTROL" > dist/$(PROG_NAME)/DEBIAN/control
	echo "$$BASED_ADVERTISEMENT" > dist/$(PROG_NAME)/DEBIAN/preinst
	dpkg-deb --build dist/$(PROG_NAME)
	cp dist/*.deb build/

build-release:
	go install .
	go build -ldflags "\
		-X 'main.ProgramName=$(PROG_NAME)' \
		-X 'main.CommitHash=$(GIT_HASH_LONG)' \
		-X 'main.Version=$(VERSION)' \
		-X 'main.BuildDate=$(BUILD_DATE)' \
		" \
		-o build/$(PROG_NAME) .
	DRCAT_BINARY_DIR=build go test -v

package-release-deb: build-release
	mkdir -p dist/$(PROG_NAME)/DEBIAN
	mkdir -p dist/$(PROG_NAME)$(DEB_INSTALL_DIR)
	cp build/$(PROG_NAME) dist/$(PROG_NAME)$(DEB_INSTALL_DIR)/$(PROG_NAME)
	mv build/$(PROG_NAME) build/$(PROG_NAME)-$(VERSION)
	touch dist/$(PROG_NAME)/DEBIAN/control
	echo "$$DEBIAN_CONTROL" > dist/$(PROG_NAME)/DEBIAN/control
	dpkg-deb --build dist/$(PROG_NAME)
	cp dist/*.deb build/$(PROG_NAME)-$(VERSION).deb
	cd build; md5sum $(PROG_NAME)-$(VERSION).deb \
		$(PROG_NAME)-$(VERSION) > $(PROG_NAME)-$(VERSION).deb.md5

sign-deb-release: package-release-deb
	@echo

test-release-deb:
	wget https://github-artifacts-zeebrow.s3.amazonaws.com/$(PROG_NAME)/releases/debian/v$(VERSION)/amd64/$(PROG_NAME)-$(VERSION).deb
	wget https://github-artifacts-zeebrow.s3.amazonaws.com/$(PROG_NAME)/releases/debian/v$(VERSION)/amd64/$(PROG_NAME)-$(VERSION).deb.md5
	md5sum -c $(PROG_NAME)-$(VERSION).deb.md5
	$(PROG_NAME) -version
	-sudo apt -y remove $(PROG_NAME)
	sudo apt install ./*.deb

release-deb: clean package-release-deb

remove-deb:
	-sudo apt -y remove $(PROG_NAME) 
reinstall-deb: clean remove-deb package-deb
	sudo apt install ./build/$(PROG_NAME).deb

clean:
	rm -rf dist/
	rm -rf build/*

.PHONY: build

