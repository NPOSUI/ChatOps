GO=go
GR=goreleaser
BUILD_FLAGS=

all: build

build:
	$(GR) release --snapshot --clean

clean:
	rm -rf dist