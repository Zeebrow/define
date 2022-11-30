PROG_NAME := define
GIT_HASH := $(shell git rev-parse --short HEAD)
GIT_HASH_LONG := $(shell git rev-parse HEAD)
BUILD_DATE := $(shell date -I)

ifndef ($(VERSION))
VERSION := 0.1.0-dev-$(GIT_HASH)
endif

TGT_FILE := cmd/define.go
GOARCH := amd64#amd64, 386, arm, ppc64
GOOS := linux#linux, darwin, windows, netbsd
DEB_INSTALL_DIR := /usr/local/bin

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

build-with-keys:
	go build -ldflags " \
		-X 'main.Version=$(VERSION)' \
		-X 'main.BuildDate=$(BUILD_DATE)' \
		-X 'main.CommitHash=$(GIT_HASH_LONG)' \
		-X 'main.ProgramName=$(PROG_NAME)' \
		-X 'main.MWDictionaryApiKey=$(MW_DICT_API_KEY)' \
		-X 'main.MWThesaurusApiKey=$(MW_THES_API_KEY)' \
		" \
		-o build/$(PROG_NAME) $(TGT_FILE) 

build:
	go build -ldflags " \
		-X 'main.Version=$(VERSION)' \
		-X 'main.BuildDate=$(BUILD_DATE)' \
		-X 'main.CommitHash=$(GIT_HASH_LONG)' \
		-X 'main.ProgramName=$(PROG_NAME)' \
		" \
		-o build/$(PROG_NAME) $(TGT_FILE) 


build-release:
	go build -ldflags " \
		-X 'main.Version=$(VERSION)' \
		-X 'main.BuildDate=$(BUILD_DATE)' \
		-X 'main.CommitHash=$(GIT_HASH_LONG)' \
		-X 'main.ProgramName=$(PROG_NAME)' \
		" \
		-o build/$(PROG_NAME)-$(VERSION) $(TGT_FILE) 

release-deb: build
	mkdir -p build/$(GOOS)/$(GOARCH)
	mkdir -p dist/$(PROG_NAME)/DEBIAN
	mkdir -p dist/$(PROG_NAME)$(DEB_INSTALL_DIR)
	cp build/$(PROG_NAME) build/$(GOOS)/$(GOARCH)/$(PROG_NAME)-$(VERSION)-$(GOOS)-$(GOARCH)
	cp build/$(PROG_NAME) dist/$(PROG_NAME)$(DEB_INSTALL_DIR)/$(PROG_NAME)
	touch dist/$(PROG_NAME)/DEBIAN/control
	echo "$$DEBIAN_CONTROL" > dist/$(PROG_NAME)/DEBIAN/control
	dpkg-deb --build dist/$(PROG_NAME)
	cp dist/*.deb build/$(PROG_NAME)-$(VERSION).deb
	cd build; md5sum \
		$(PROG_NAME)-$(VERSION).deb \
		$(GOOS)/$(GOARCH)/$(PROG_NAME)-$(VERSION)-$(GOOS)-$(GOARCH) \
	> SUMS.md5
	cd build; md5sum -c SUMS.md5

package-deb: build
	mkdir -p build/$(GOOS)/$(GOARCH)
	mkdir -p dist/$(PROG_NAME)/DEBIAN
	mkdir -p dist/$(PROG_NAME)$(DEB_INSTALL_DIR)
	cp build/$(PROG_NAME) build/$(GOOS)/$(GOARCH)/$(PROG_NAME)-$(VERSION)-$(GOOS)-$(GOARCH)
	cp build/$(PROG_NAME) dist/$(PROG_NAME)$(DEB_INSTALL_DIR)/$(PROG_NAME)
	touch dist/$(PROG_NAME)/DEBIAN/control
	echo "$$DEBIAN_CONTROL" > dist/$(PROG_NAME)/DEBIAN/control
	dpkg-deb --build dist/$(PROG_NAME)
	cp dist/*.deb build/$(PROG_NAME)-$(VERSION).deb
	cd build; md5sum \
		$(PROG_NAME)-$(VERSION).deb \
		$(GOOS)/$(GOARCH)/$(PROG_NAME)-$(VERSION)-$(GOOS)-$(GOARCH) \
	> SUMS.md5
	cd build; md5sum -c SUMS.md5

clean:
	rm -rf dist/
	rm -rf build/*

reinstall-deb: clean package-deb
	sudo apt install ./build/*.deb
reinstall-deb-release: clean release-deb
	sudo apt install ./build/*.deb


.PHONY: build test
