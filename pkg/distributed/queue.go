package distributed

import (
    "context"
    "fmt"
    "sync"
    "time"
)

// TaskQueue interface for distributed task queue
type TaskQueue interface {
    Enqueue(task *Task) error
    Dequeue() (*Task, error)
    Ack(taskID string) error
    Retry(task *Task) error
    Size() int
    Close() error
}

// MemoryQueue implements in-memory queue (for development)
type MemoryQueue struct {
    queue      chan *Task
    pending    map[string]*Task
    mu         sync.RWMutex
    maxSize    int
    ctx        context.Context
    cancel     context.CancelFunc
}

// NewMemoryQueue creates a new in-memory queue
func NewMemoryQueue(maxSize int) *MemoryQueue {
    ctx, cancel := context.WithCancel(context.Background())
    return &MemoryQueue{
        queue:   make(chan *Task, maxSize),
        pending: make(map[string]*Task),
        maxSize: maxSize,
        ctx:     ctx,
        cancel:  cancel,
    }
}

// Enqueue adds a task to the queue
func (q *MemoryQueue) Enqueue(task *Task) error {
    select {
    case q.queue <- task:
        q.mu.Lock()
        q.pending[task.ID] = task
        q.mu.Unlock()
        return nil
    case <-time.After(5 * time.Second):
        return fmt.Errorf("queue full, task %s not enqueued", task.ID)
    }
}

// Dequeue retrieves a task from the queue
func (q *MemoryQueue) Dequeue() (*Task, error) {
    select {
    case task := <-q.queue:
        task.Status = StatusRunning
        task.UpdatedAt = time.Now()
        return task, nil
    case <-q.ctx.Done():
        return nil, fmt.Errorf("queue closed")
    }
}

// Ack acknowledges task completion
func (q *MemoryQueue) Ack(taskID string) error {
    q.mu.Lock()
    defer q.mu.Unlock()
    delete(q.pending, taskID)
    return nil
}

// Retry requeues a failed task
func (q *MemoryQueue) Retry(task *Task) error {
    task.Status = StatusPending
    task.UpdatedAt = time.Now()
    return q.Enqueue(task)
}

// Size returns current queue size
func (q *MemoryQueue) Size() int {
    return len(q.queue)
}

// Close closes the queue
func (q *MemoryQueue) Close() error {
    q.cancel()
    close(q.queue)
    return nil
}

// QueueStats provides queue statistics
type QueueStats struct {
    TotalTasks    int `json:"total_tasks"`
    PendingTasks  int `json:"pending_tasks"`
    RunningTasks  int `json:"running_tasks"`
    CompletedTasks int `json:"completed_tasks"`
    FailedTasks   int `json:"failed_tasks"`
}

// GetStats returns queue statistics
func (q *MemoryQueue) GetStats() *QueueStats {
    q.mu.RLock()
    defer q.mu.RUnlock()
    
    stats := &QueueStats{
        TotalTasks:   len(q.pending),
        PendingTasks: len(q.queue),
    }
    
    for _, task := range q.pending {
        switch task.Status {
        case StatusRunning:
            stats.RunningTasks++
        case StatusCompleted:
            stats.CompletedTasks++
        case StatusFailed:
            stats.FailedTasks++
        }
    }
    
    return stats
}
