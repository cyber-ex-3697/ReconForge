.PHONY: build clean test install docker run help

VERSION := $(shell git describe --tags 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u +%Y%m%d.%H%M%S)
LDFLAGS := -ldflags "-X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME)"

build:
	go build $(LDFLAGS) -o reconforge cmd/reconforge/main.go

clean:
	rm -rf reconforge output/ logs/

test:
	go test -v ./...

install:
	go build -o /usr/local/bin/reconforge cmd/reconforge/main.go

docker:
	docker build -t reconforge:latest .

docker-run:
	docker run --rm reconforge:latest -t example.com

run:
	go run cmd/reconforge/main.go -t example.com

help:
	@echo "Available targets:"
	@echo "  build     - Build binary"
	@echo "  clean     - Clean build artifacts"
	@echo "  test      - Run tests"
	@echo "  install   - Install to /usr/local/bin"
	@echo "  docker    - Build Docker image"
	@echo "  run       - Run basic scan"
	@echo "  help      - Show this help"

.DEFAULT_GOAL := help
