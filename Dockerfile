FROM golang:1.21-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o reconforge cmd/reconforge/main.go

FROM alpine:3.19

RUN apk add --no-cache \
    ca-certificates \
    bash \
    curl \
    wget \
    nmap \
    jq \
    parallel \
    && rm -rf /var/cache/apk/*

RUN wget -q https://github.com/projectdiscovery/subfinder/releases/latest/download/subfinder_linux_amd64.zip && \
    unzip subfinder_linux_amd64.zip && mv subfinder /usr/local/bin/ && rm subfinder_linux_amd64.zip

RUN wget -q https://github.com/projectdiscovery/httpx/releases/latest/download/httpx_linux_amd64.zip && \
    unzip httpx_linux_amd64.zip && mv httpx /usr/local/bin/ && rm httpx_linux_amd64.zip

RUN wget -q https://github.com/projectdiscovery/nuclei/releases/latest/download/nuclei_linux_amd64.zip && \
    unzip nuclei_linux_amd64.zip && mv nuclei /usr/local/bin/ && rm nuclei_linux_amd64.zip

RUN go install github.com/tomnomnom/assetfinder@latest && cp /root/go/bin/assetfinder /usr/local/bin/

WORKDIR /app

COPY --from=builder /build/reconforge .
COPY config.yaml .

RUN mkdir -p /app/output

ENTRYPOINT ["./reconforge"]
CMD ["-h"]
