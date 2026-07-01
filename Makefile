VERSION ?= dev

.PHONY: build build-all clean

build:
	go build -ldflags="-X github.com/xdewx/git-ai-exporter/cmd.Version=$(VERSION)" -o dist/git-ai-exporter .

build-all:
	GOOS=linux   GOARCH=amd64 go build -ldflags="-X github.com/xdewx/git-ai-exporter/cmd.Version=$(VERSION)" -o dist/git-ai-exporter-linux-amd64 .
	GOOS=darwin  GOARCH=amd64 go build -ldflags="-X github.com/xdewx/git-ai-exporter/cmd.Version=$(VERSION)" -o dist/git-ai-exporter-darwin-amd64 .
	GOOS=darwin  GOARCH=arm64 go build -ldflags="-X github.com/xdewx/git-ai-exporter/cmd.Version=$(VERSION)" -o dist/git-ai-exporter-darwin-arm64 .
	GOOS=windows GOARCH=amd64 go build -ldflags="-X github.com/xdewx/git-ai-exporter/cmd.Version=$(VERSION)" -o dist/git-ai-exporter-windows-amd64.exe .

clean:
	rm -rf dist/
