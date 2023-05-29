BUILD_FLAGS=

all: build

build:
	goreleaser release --snapshot --clean

build_auth:
	go build -o executor_auth_linux_amd64 botkube/external-plugins/executors/auth/main.go

clean:
	rm -rf dist

clean_tmp:
	rm -rf /tmp/*