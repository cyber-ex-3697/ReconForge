package passive_recon

import (
    "sync"
    "time"
)

// PassiveRateLimiter manages rate limiting for passive API calls
type PassiveRateLimiter struct {
    limits    map[string]*RateLimit
    defaultLimit int
    mu        sync.RWMutex
}

// RateLimit represents rate limit for a specific API
type RateLimit struct {
    RequestsPerSecond int
    LastRequest       time.Time
    Tokens            int
    mu                sync.Mutex
}

// NewPassiveRateLimiter creates a new rate limiter
func NewPassiveRateLimiter(defaultRPS int) *PassiveRateLimiter {
    return &PassiveRateLimiter{
        limits:       make(map[string]*RateLimit),
        defaultLimit: defaultRPS,
    }
}

// SetLimit sets rate limit for an API
func (prl *PassiveRateLimiter) SetLimit(api string, rps int) {
    prl.mu.Lock()
    defer prl.mu.Unlock()
    
    prl.limits[api] = &RateLimit{
        RequestsPerSecond: rps,
        Tokens:            rps,
        LastRequest:       time.Now(),
    }
}

// Wait waits for rate limit to allow request
func (prl *PassiveRateLimiter) Wait(api string) {
    prl.mu.RLock()
    limit, exists := prl.limits[api]
    prl.mu.RUnlock()
    
    if !exists {
        limit = &RateLimit{
            RequestsPerSecond: prl.defaultLimit,
            Tokens:            prl.defaultLimit,
            LastRequest:       time.Now(),
        }
    }
    
    limit.mu.Lock()
    defer limit.mu.Unlock()
    
    now := time.Now()
    elapsed := now.Sub(limit.LastRequest)
    
    // Refill tokens based on elapsed time
    refill := int(elapsed.Seconds() * float64(limit.RequestsPerSecond))
    if refill > 0 {
        limit.Tokens += refill
        if limit.Tokens > limit.RequestsPerSecond {
            limit.Tokens = limit.RequestsPerSecond
        }
        limit.LastRequest = now
    }
    
    if limit.Tokens <= 0 {
        waitTime := time.Second / time.Duration(limit.RequestsPerSecond)
        time.Sleep(waitTime)
        limit.Tokens = limit.RequestsPerSecond - 1
        limit.LastRequest = time.Now()
    } else {
        limit.Tokens--
    }
}

// GetRemainingTokens returns remaining tokens for an API
func (prl *PassiveRateLimiter) GetRemainingTokens(api string) int {
    prl.mu.RLock()
    limit, exists := prl.limits[api]
    prl.mu.RUnlock()
    
    if !exists {
        return prl.defaultLimit
    }
    
    limit.mu.Lock()
    defer limit.mu.Unlock()
    return limit.Tokens
}

// Reset resets rate limit for an API
func (prl *PassiveRateLimiter) Reset(api string) {
    prl.mu.Lock()
    defer prl.mu.Unlock()
    
    if limit, exists := prl.limits[api]; exists {
        limit.Tokens = limit.RequestsPerSecond
        limit.LastRequest = time.Now()
    }
}

// AddBackoff adds exponential backoff for rate-limited APIs
func (prl *PassiveRateLimiter) AddBackoff(api string, attempt int) {
    backoff := time.Duration(1<<attempt) * time.Second
    if backoff > 30*time.Second {
        backoff = 30 * time.Second
    }
    time.Sleep(backoff)
}

// CreateAPILimiter creates a limiter for a specific API
func (prl *PassiveRateLimiter) CreateAPILimiter(api string, rps int) func() {
    prl.SetLimit(api, rps)
    return func() {
        prl.Wait(api)
    }
}

// GetDefaultLimit returns default rate limit
func (prl *PassiveRateLimiter) GetDefaultLimit() int {
    return prl.defaultLimit
}

// SetDefaultLimit sets default rate limit
func (prl *PassiveRateLimiter) SetDefaultLimit(rps int) {
    prl.defaultLimit = rps
}
