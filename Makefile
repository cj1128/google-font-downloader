.PHONY: install build-osx build-linux build-windows build-all

install:
	godep go install -ldflags "-X main.appVersion=`cat VERSION`"

build-all: build-osx build-linux build-windows

build-osx:
	godep go build -ldflags "-s -w -X main.appVersion=`cat VERSION`" -o tmp/google-font-downloader-osx

build-linux:
	GOOS=linux godep go build -o tmp/google-font-downloader-linux -ldflags "-s -w -X main.appVersion=`cat VERSION`"

build-windows:
	GOOS=windows godep go build -o tmp/google-font-downloader-windows -ldflags "-s -w -X main.appVersion=`cat VERSION`"
