BIN_NAME := watchdog
BIN_PATH := bin/$(BIN_NAME)
BUILD_FLAGS := -installsuffix "static"

.PHONY: all build clean deps

all: build

build:
	CGO_ENABLED=1 \
	GO111MODULE=on \
	GOARCH=amd64 \
	go build \
	$(BUILD_FLAGS) \
	-o $(BIN_PATH) \
	$$(find cmd/app/*.go)

clean:
	@echo Cleaning bin/ directory... && \
		rm -rfv bin/

deps:
	@echo Downloading go.mod dependencies && \
		go mod download

install: build
	@echo Installing...
	@cp -av bin/watchdog /usr/bin/watchdog
	@chmod +x /usr/bin/watchdog
	@mkdir -p /etc/watchdog
	@cp -av dist/config.example.json /etc/watchdog/config.json
	@cp -av dist/watchdog.service /etc/systemd/system
	@chmod 644 /etc/systemd/system/watchdog.service