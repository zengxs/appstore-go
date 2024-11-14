SOURCES := $(shell find . -name '*.go')

.PHONY: all

all: dist/macos/appstore dist/linux-amd64/appstore dist/linux-arm64/appstore dist/windows-amd64/appstore.exe

.PHONY: clean

clean:
	rm -rf dist

dist/macos/appstore: $(SOURCES)
	GOOS=darwin GOARCH=amd64 go build -o dist/macos/appstore-amd64 ./cmd/appstore
	GOOS=darwin GOARCH=arm64 go build -o dist/macos/appstore-arm64 ./cmd/appstore
	go run ./cmd/lipo -output $@ -create dist/macos/appstore-amd64 dist/macos/appstore-arm64
	rm dist/macos/appstore-amd64 dist/macos/appstore-arm64

dist/linux-amd64/appstore: $(SOURCES)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $@ ./cmd/appstore

dist/linux-arm64/appstore: $(SOURCES)
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o $@ ./cmd/appstore

dist/windows-amd64/appstore.exe: $(SOURCES)
	GOOS=windows GOARCH=amd64 go build -o $@ ./cmd/appstore
