# ReconForge - Docker Deployment Guide

## Quick Start

```bash
# Pull image
docker pull reconforge/reconforge:latest

# Run basic scan
docker run --rm reconforge/reconforge:latest -t example.com

# Run deep scan
docker run --rm reconforge/reconforge:latest -t example.com --deep

Docker Compose

docker-compose.yml

version: '3.8'

services:
  reconforge:
    image: reconforge/reconforge:latest
    container_name: reconforge
    volumes:
      - ./output:/app/output
      - ./config.yaml:/app/config.yaml
    environment:
      - TARGET=example.com
      - DEEP_SCAN=false
      - THREADS=50
    command: -t ${TARGET} ${DEEP_SCAN:+-d} -T ${THREADS}

Run with Compose

export TARGET=example.com
export DEEP_SCAN=true
docker-compose up


Build Custom Image

Dockerfile

FROM reconforge/reconforge:latest AS builder

COPY custom_plugins/ /app/plugins/community/

ENTRYPOINT ["./reconforge"]

Build

docker build -t reconforge-custom .

Volume Mounts

Output Directory

docker run --rm -v $(pwd)/output:/app/output reconforge/reconforge:latest -t example.com

Custom Config

docker run --rm -v $(pwd)/config.yaml:/app/config.yaml reconforge/reconforge:latest -t example.com

Custom Wordlists

docker run --rm -v $(pwd)/wordlists:/app/wordlists reconforge/reconforge:latest -t example.com


Environment Variables


Variable		Description		Default
TARGET			Target domain		(required)
DEEP_SCAN		Enable deep scan	false
THREADS			Number of threads	50
OUTPUT_DIR		Output directory	/app/output


Kubernetes Deployment

deployment.yaml

apiVersion: apps/v1
kind: Deployment
metadata:
  name: reconforge
spec:
  replicas: 1
  selector:
    matchLabels:
      app: reconforge
  template:
    metadata:
      labels:
        app: reconforge
    spec:
      containers:
      - name: reconforge
        image: reconforge/reconforge:latest
        args: ["-t", "example.com", "--deep"]
        volumeMounts:
        - name: output
          mountPath: /app/output
      volumes:
      - name: output
        persistentVolumeClaim:
          claimName: reconforge-pvc

CI/CD Integration

GitHub Actions

name: ReconForge Scan
on:
  schedule:
    - cron: '0 0 * * *'
jobs:
  scan:
    runs-on: ubuntu-latest
    steps:
      - name: Run ReconForge
        run: |
          docker run --rm reconforge/reconforge:latest -t example.com

Multiple Instances

Distributed Scanning

version: '3.8'
services:
  scan1:
    image: reconforge/reconforge:latest
    command: -t sub1.example.com
  scan2:
    image: reconforge/reconforge:latest
    command: -t sub2.example.com

Resource Limits

docker run --rm \
  --memory="4g" \
  --cpus="2" \
  reconforge/reconforge:latest -t example.com

Docker Hub

Official images: https://hub.docker.com/r/reconforge/reconforge


Troubleshooting

Permission Denied

# Fix volume permissions
sudo chown -R 1000:1000 ./output

Out of Memory

# Increase memory limit
docker run --rm --memory="8g" reconforge/reconforge:latest -t example.com


