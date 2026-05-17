# Distributed Scan Example

## Overview

Distributed scanning allows running ReconForge across multiple nodes for large targets.

## Architecture


┌─────────────┐
│ Master │
│ Node │
└──────┬──────┘
│
┌────────┼────────┐
│ │ │
▼ ▼ ▼
┌─────────┐┌─────────┐┌─────────┐
│ Worker1 ││ Worker2 ││ Worker3 │
└─────────┘└─────────┘└─────────┘


## Prerequisites

- Redis server (for coordination)
- All nodes have ReconForge installed
- Network connectivity between nodes

## Setup Redis

```bash
docker run -d --name redis -p 6379:6379 redis:alpine

Master Node Configuration

config-master.yaml

distributed:
  enabled: true
  role: master
  redis: "localhost:6379"
  workers: 3

Run Master

./reconforge -t example.com --distributed --master

Worker Node Configuration

config-worker.yaml

distributed:
  enabled: true
  role: worker
  redis: "master-ip:6379"

Run Worker

./reconforge --worker

Using Docker Compose

docker-compose-distributed.yml

version: '3.8'

services:
  redis:
    image: redis:alpine
    container_name: redis
    ports:
      - "6379:6379"

  master:
    image: reconforge/reconforge:latest
    container_name: master
    depends_on:
      - redis
    command: -t example.com --deep --distributed --master
    environment:
      - REDIS_ADDR=redis:6379

  worker1:
    image: reconforge/reconforge:latest
    container_name: worker1
    depends_on:
      - redis
    command: --worker
    environment:
      - REDIS_ADDR=redis:6379

  worker2:
    image: reconforge/reconforge:latest
    container_name: worker2
    depends_on:
      - redis
    command: --worker
    environment:
      - REDIS_ADDR=redis:6379

  worker3:
    image: reconforge/reconforge:latest
    container_name: worker3
    depends_on:
      - redis
    command: --worker
    environment:
      - REDIS_ADDR=redis:6379

Run

docker-compose -f docker-compose-distributed.yml up


Kubernetes Distributed Scan

job.yaml

apiVersion: batch/v1
kind: Job
metadata:
  name: reconforge-distributed
spec:
  parallelism: 3
  completions: 3
  template:
    spec:
      containers:
      - name: reconforge
        image: reconforge/reconforge:latest
        args: ["-t", "example.com", "--deep", "--distributed"]
      restartPolicy: Never


Performance Comparison

Configuration		Time		Nodes

Single Node		30 min		1
3 Nodes			10 min		3
5 Nodes			06 min		5


Best Practices

Use same version on all nodes

Monitor Redis for queue size

Balance workload across nodes

Collect results from master only

Handle failures with retries


Troubleshooting

Node Not Connecting

# Check Redis connectivity
redis-cli -h master-ip ping

Uneven Work Distribution

# Check queue size
redis-cli LLEN reconforge:queue

Partial Results

# Resume from checkpoint on master
./reconforge -r checkpoint.json


