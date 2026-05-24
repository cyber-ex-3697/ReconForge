.PHONY: build clean lint test fmt

build:
	go build -o reconforge cmd/reconforge/main.go

clean:
	rm -f reconforge
	rm -rf recon_*/

lint:
	golangci-lint run

fmt:
	go fmt ./...

test:
	go test -v ./...

run:
	./reconforge -t example.com

install:
	go install ./cmd/reconforge

.PHONY: all
all: fmt lint test build
