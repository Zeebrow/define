GIT_HASH := $(shell git rev-parse --short HEAD)
GIT_HASH_LONG := $(shell git rev-parse HEAD)
GOARCH := amd64 #amd64, 386, arm, ppc64
GOOS := linux #linux, darwin, windows, netbsd

build-linux-amd64:
	go install .
	GOOS=linux GOARCH=amd64 \
			 go build -ldflags "-X 'main.Version=$(GIT_HASH_LONG)'" -o build/define-$(GIT_HASH)-linux-amd64 .

build-windows-amd64:
	go install .
	GOOS=windows GOARCH=arm \
			 go build -ldflags "-X 'main.Version=$(GIT_HASH_LONG)'" -o build/define-$(GIT_HASH)-windows-amd64 .

build:
	go install .
	go build -ldflags "-X 'main.Version=$(GIT_HASH_LONG)'" -o build/define-$(GIT_HASH) .

clean:
	rm -rf build/*

.PHONY: build clean
