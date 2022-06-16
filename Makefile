PROG_NAME := define
VERSION := 0.1
GIT_HASH := $(shell git rev-parse --short HEAD)
GIT_HASH_LONG := $(shell git rev-parse HEAD)
BUILD_DATE := $(shell date -I)
GOARCH := amd64 #amd64, 386, arm, ppc64
GOOS := linux #linux, darwin, windows, netbsd
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
Homepage: https://github.com/zeebrow/define
Description: Use the Merriam Webster API to find definitions for words from the terminal
endef
export DEBIAN_CONTROL

install-zeebrow:
	go install .
	go build \
		-ldflags " \
		-X 'main.Version=$(GIT_HASH_LONG)' \
		-X 'main.BuildDate=$(BUILD_DATE)' \
		" \
		-o build/$(PROG_NAME) .
	cp build/$(PROG_NAME) $(HOME)/.local/bin/scripts

build-dev:
	go install .
	go build -ldflags " \
		-X 'main.Version=$(GIT_HASH_LONG)' \
		-X 'main.BuildDate=$(BUILD_DATE)' \
		" \
		-o build/define-$(GIT_HASH) .

build:
	go install .
	go build -ldflags " \
		-X 'main.Version=dev-$(GIT_HASH)' \
		-X 'main.BuildDate=$(BUILD_DATE)' \
		-X 'main.CommitHash=$(GIT_HASH_LONG)' \
		-X 'main.ProgramName=$(PROG_NAME)' \
		-X 'main.ProgInfo.Test=WHOOPIE' \
		" \
		-o build/define .

package-deb: build
	mkdir -p dist/$(PROG_NAME)/DEBIAN
	mkdir -p dist/$(PROG_NAME)$(DEB_INSTALL_DIR)
	cp build/$(PROG_NAME) dist/$(PROG_NAME)$(DEB_INSTALL_DIR)/$(PROG_NAME)
	touch dist/$(PROG_NAME)/DEBIAN/control
	echo "$$DEBIAN_CONTROL" > dist/$(PROG_NAME)/DEBIAN/control
	dpkg-deb --build dist/$(PROG_NAME)
	cp dist/*.deb build/

clean:
	rm -rf dist/
	rm -rf build/*

.PHONY: build clean
