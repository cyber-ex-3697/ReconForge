package scheduler

import (
    "sync"
    "time"
)

type AdaptiveWorkerPool struct {
    pool           *WorkerPool
    currentWorkers int
    maxWorkers     int
    minWorkers     int
    rateLimitHits  int
    mu             sync.Mutex
}

func NewAdaptiveWorkerPool(minWorkers, maxWorkers int) *AdaptiveWorkerPool {
    pool := NewWorkerPool(minWorkers)
    pool.Start()
    
    return &AdaptiveWorkerPool{
        pool:           pool,
        currentWorkers: minWorkers,
        maxWorkers:     maxWorkers,
        minWorkers:     minWorkers,
        rateLimitHits:  0,
    }
}

func (p *AdaptiveWorkerPool) RecordRateLimit() {
    p.mu.Lock()
    defer p.mu.Unlock()
    p.rateLimitHits++
    
    if p.rateLimitHits > 5 {
        newWorkers := p.currentWorkers / 2
        if newWorkers < p.minWorkers {
            newWorkers = p.minWorkers
        }
        p.currentWorkers = newWorkers
        p.rateLimitHits = 0
    }
}

func (p *AdaptiveWorkerPool) RecordSuccess() {
    p.mu.Lock()
    defer p.mu.Unlock()
    
    if p.rateLimitHits > 0 {
        p.rateLimitHits--
    } else if p.currentWorkers < p.maxWorkers {
        newWorkers := p.currentWorkers * 2
        if newWorkers > p.maxWorkers {
            newWorkers = p.maxWorkers
        }
        p.currentWorkers = newWorkers
    }
}

func (p *AdaptiveWorkerPool) GetCurrentWorkers() int {
    p.mu.Lock()
    defer p.mu.Unlock()
    return p.currentWorkers
}

func (p *AdaptiveWorkerPool) Submit(job Job) {
    p.pool.Submit(job)
}

func (p *AdaptiveWorkerPool) Results() <-chan Result {
    return p.pool.Results()
}

func (p *AdaptiveWorkerPool) Stop() {
    p.pool.Stop()
}
