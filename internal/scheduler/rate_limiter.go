package scheduler

import (
    "sync"
    "time"
)

type RateLimiter struct {
    rate       int           // requests per second
    interval   time.Duration
    tokens     int
    maxTokens  int
    mu         sync.Mutex
    lastRefill time.Time
}

func NewRateLimiter(rate int) *RateLimiter {
    return &RateLimiter{
        rate:       rate,
        interval:   time.Second / time.Duration(rate),
        tokens:     rate,
        maxTokens:  rate,
        lastRefill: time.Now(),
    }
}

func (r *RateLimiter) Allow() bool {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    now := time.Now()
    elapsed := now.Sub(r.lastRefill)
    
    // Refill tokens
    if elapsed >= r.interval {
        newTokens := int(elapsed / r.interval)
        r.tokens += newTokens
        if r.tokens > r.maxTokens {
            r.tokens = r.maxTokens
        }
        r.lastRefill = now
    }
    
    if r.tokens > 0 {
        r.tokens--
        return true
    }
    return false
}

func (r *RateLimiter) Wait() {
    for !r.Allow() {
        time.Sleep(10 * time.Millisecond)
    }
}

func (r *RateLimiter) SetRate(rate int) {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.rate = rate
    r.interval = time.Second / time.Duration(rate)
    r.maxTokens = rate
}

func (r *RateLimiter) GetRate() int {
    r.mu.Lock()
    defer r.mu.Unlock()
    return r.rate
}
