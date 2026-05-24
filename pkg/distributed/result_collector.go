package distributed

import (
    "sync"
    "time"
)

// ResultCollector collects and aggregates task results
type ResultCollector struct {
    results     map[string]*TaskResult
    pending     map[string]bool
    mu          sync.RWMutex
    callbacks   []ResultCallback
}

// ResultCallback is called when a result is received
type ResultCallback func(taskID string, result *TaskResult)

// NewResultCollector creates a new result collector
func NewResultCollector() *ResultCollector {
    return &ResultCollector{
        results:   make(map[string]*TaskResult),
        pending:   make(map[string]bool),
        callbacks: make([]ResultCallback, 0),
    }
}

// AddResult adds a task result
func (c *ResultCollector) AddResult(taskID string, result *TaskResult) {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    c.results[taskID] = result
    delete(c.pending, taskID)
    
    // Notify callbacks
    for _, cb := range c.callbacks {
        cb(taskID, result)
    }
}

// GetResult retrieves a result
func (c *ResultCollector) GetResult(taskID string) *TaskResult {
    c.mu.RLock()
    defer c.mu.RUnlock()
    return c.results[taskID]
}

// RegisterPending registers a pending task
func (c *ResultCollector) RegisterPending(taskID string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.pending[taskID] = true
}

// GetPending returns all pending task IDs
func (c *ResultCollector) GetPending() []string {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    var pending []string
    for id := range c.pending {
        pending = append(pending, id)
    }
    return pending
}

// GetCompleted returns all completed task IDs
func (c *ResultCollector) GetCompleted() []string {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    var completed []string
    for id := range c.results {
        completed = append(completed, id)
    }
    return completed
}

// AddCallback adds a result callback
func (c *ResultCollector) AddCallback(cb ResultCallback) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.callbacks = append(c.callbacks, cb)
}

// WaitForAll waits for all pending tasks to complete
func (c *ResultCollector) WaitForAll(timeout time.Duration) bool {
    deadline := time.Now().Add(timeout)
    
    for {
        if time.Now().After(deadline) {
            return false
        }
        
        c.mu.RLock()
        pendingCount := len(c.pending)
        c.mu.RUnlock()
        
        if pendingCount == 0 {
            return true
        }
        
        time.Sleep(1 * time.Second)
    }
}

// AggregateResults aggregates results by task type
func (c *ResultCollector) AggregateResults() map[TaskType]interface{} {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    aggregated := make(map[TaskType]interface{})
    
    for _, result := range c.results {
        if result.Success && result.Data != nil {
            // This would need proper type assertion in production
            aggregated[TaskType("unknown")] = result.Data
        }
    }
    
    return aggregated
}

// Clear clears all collected results
func (c *ResultCollector) Clear() {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    c.results = make(map[string]*TaskResult)
    c.pending = make(map[string]bool)
}
