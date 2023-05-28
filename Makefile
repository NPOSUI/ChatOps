BUILD_FLAGS=

all: build

build:
	goreleaser release --snapshot --clean

clean:
	rm -rf dist