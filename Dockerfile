FROM golang:1.21-alpine AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o reconforge cmd/reconforge/main.go

FROM alpine:3.19
RUN apk add --no-cache ca-certificates curl wget nmap jq
COPY --from=builder /build/reconforge /usr/local/bin/
ENTRYPOINT ["reconforge"]
CMD ["-h"]
