# compile commands
GO=go
GR=goreleaser
# compile option
BUILD_FLAGS=

# default
all: build

# Compile
build:
	$(GR) build release --snapshot --clean

# run
#run: build
#	./

# delete the generated file
clean:
	rm -rf dist