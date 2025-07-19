PROG_NAME := define
GIT_HASH := $(shell git rev-parse --short HEAD)
GIT_HASH_LONG := $(shell git rev-parse HEAD)
BUILD_DATE := $(shell date -I)

TGT_FILE := cmd/define.go
GOARCH := amd64#amd64, 386, arm, ppc64
GOOS := linux#linux, darwin, windows, netbsd
DEB_INSTALL_DIR := /usr/local/bin


require_version:
ifeq ($(origin VERSION),undefined)
	$(error VERSION is unset and required)
endif

deb_ctrl_file: require_version
	$(shell mkdir -p dist/$(PROG_NAME)/DEBIAN)
	$(file > dist/$(PROG_NAME)/DEBIAN/control,Package: $(PROG_NAME))
	$(file >> dist/$(PROG_NAME)/DEBIAN/control,Version: $(VERSION))
	$(file >> dist/$(PROG_NAME)/DEBIAN/control,Provides: $(PROG_NAME))
	$(file >> dist/$(PROG_NAME)/DEBIAN/control,Section: custom)
	$(file >> dist/$(PROG_NAME)/DEBIAN/control,Priority: optional)
	$(file >> dist/$(PROG_NAME)/DEBIAN/control,Architecture: $(GOARCH))
	$(file >> dist/$(PROG_NAME)/DEBIAN/control,Essential: no)
	$(file >> dist/$(PROG_NAME)/DEBIAN/control,Installed-Size: 8192)
	$(file >> dist/$(PROG_NAME)/DEBIAN/control,Maintainer: zeebrow)
	$(file >> dist/$(PROG_NAME)/DEBIAN/control,Homepage: https://github.com/zeebrow/define)
	$(file >> dist/$(PROG_NAME)/DEBIAN/control,Description: Use the Merriam Webster API to find definitions for words from the terminal)


build-with-keys:
	go build -ldflags " \
		-X 'github.com/Zeebrow/define/define.Version=$(VERSION)' \
		-X 'github.com/Zeebrow/define/define.BuildDate=$(BUILD_DATE)' \
		-X 'github.com/Zeebrow/define/define.CommitHash=$(GIT_HASH_LONG)' \
		-X 'github.com/Zeebrow/define/define.ProgramName=$(PROG_NAME)' \
		-X 'github.com/Zeebrow/define/define.MWDictionaryApiKey=$(MW_DICT_API_KEY)' \
		-X 'github.com/Zeebrow/define/define.MWThesaurusApiKey=$(MW_THES_API_KEY)' \
		" \
		-o build/$(PROG_NAME) $(TGT_FILE) 

build:
	go build -ldflags " \
		-X 'github.com/Zeebrow/define/define.Version=$(GIT_HASH)' \
		-X 'github.com/Zeebrow/define/define.BuildDate=$(BUILD_DATE)' \
		-X 'github.com/Zeebrow/define/define.CommitHash=$(GIT_HASH_LONG)' \
		-X 'github.com/Zeebrow/define/define.ProgramName=$(PROG_NAME)' \
		" \
		-o build/$(PROG_NAME) $(TGT_FILE) 


build-release: require_version
	go build -ldflags " \
		-X 'github.com/Zeebrow/define/define.Version=$(VERSION)' \
		-X 'github.com/Zeebrow/define/define.BuildDate=$(BUILD_DATE)' \
		-X 'github.com/Zeebrow/define/define.CommitHash=$(GIT_HASH_LONG)' \
		-X 'github.com/Zeebrow/define/define.ProgramName=$(PROG_NAME)' \
		" \
		-o build/$(PROG_NAME)-$(VERSION) $(TGT_FILE) 

package-deb: deb_ctrl_file build-release
	mkdir -p build/$(GOOS)/$(GOARCH)
	mkdir -p dist/$(PROG_NAME)/DEBIAN
	mkdir -p dist/$(PROG_NAME)$(DEB_INSTALL_DIR)
	cp build/$(PROG_NAME) build/$(GOOS)/$(GOARCH)/$(PROG_NAME)-$(VERSION)-$(GOOS)-$(GOARCH)
	cp build/$(PROG_NAME) dist/$(PROG_NAME)$(DEB_INSTALL_DIR)/$(PROG_NAME)
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


.PHONY: build test
