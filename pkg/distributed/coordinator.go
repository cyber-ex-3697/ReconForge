package distributed

import (
    "fmt"
    "sync"
    "time"
)

// Coordinator manages distributed scan across nodes
type Coordinator struct {
    nodeID      string
    queue       TaskQueue
    workers     map[string]*WorkerInfo
    tasks       map[string]*Task
    results     map[string]*TaskResult
    mu          sync.RWMutex
    isMaster    bool
}

// WorkerInfo holds worker metadata
type WorkerInfo struct {
    ID        string
    Status    WorkerStatus
    LastSeen  time.Time
    TasksDone int
    TasksFailed int
}

// NewCoordinator creates a new coordinator
func NewCoordinator(nodeID string, isMaster bool) *Coordinator {
    return &Coordinator{
        nodeID:   nodeID,
        workers:  make(map[string]*WorkerInfo),
        tasks:    make(map[string]*Task),
        results:  make(map[string]*TaskResult),
        isMaster: isMaster,
    }
}

// SetQueue sets the task queue
func (c *Coordinator) SetQueue(queue TaskQueue) {
    c.queue = queue
}

// RegisterWorker registers a worker
func (c *Coordinator) RegisterWorker(workerID string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    c.workers[workerID] = &WorkerInfo{
        ID:        workerID,
        Status:    WorkerIdle,
        LastSeen:  time.Now(),
        TasksDone: 0,
        TasksFailed: 0,
    }
}

// UnregisterWorker removes a worker
func (c *Coordinator) UnregisterWorker(workerID string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    delete(c.workers, workerID)
}

// SubmitTask submits a task to the queue
func (c *Coordinator) SubmitTask(task *Task) error {
    c.mu.Lock()
    c.tasks[task.ID] = task
    c.mu.Unlock()
    
    return c.queue.Enqueue(task)
}

// SubmitBatch submits multiple tasks
func (c *Coordinator) SubmitBatch(tasks []*Task) error {
    for _, task := range tasks {
        if err := c.SubmitTask(task); err != nil {
            return err
        }
    }
    return nil
}

// GetTaskStatus returns task status
func (c *Coordinator) GetTaskStatus(taskID string) TaskStatus {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    if task, ok := c.tasks[taskID]; ok {
        return task.Status
    }
    return StatusFailed
}

// GetResult returns task result
func (c *Coordinator) GetResult(taskID string) *TaskResult {
    c.mu.RLock()
    defer c.mu.RUnlock()
    return c.results[taskID]
}

// UpdateWorkerHeartbeat updates worker's last seen time
func (c *Coordinator) UpdateWorkerHeartbeat(workerID string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    if worker, ok := c.workers[workerID]; ok {
        worker.LastSeen = time.Now()
    }
}

// GetWorkerStats returns worker statistics
func (c *Coordinator) GetWorkerStats() map[string]*WorkerInfo {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    result := make(map[string]*WorkerInfo)
    for k, v := range c.workers {
        result[k] = &WorkerInfo{
            ID:         v.ID,
            Status:     v.Status,
            LastSeen:   v.LastSeen,
            TasksDone:  v.TasksDone,
            TasksFailed: v.TasksFailed,
        }
    }
    return result
}

// GetTaskStats returns task statistics
func (c *Coordinator) GetTaskStats() map[TaskStatus]int {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    stats := make(map[TaskStatus]int)
    for _, task := range c.tasks {
        stats[task.Status]++
    }
    return stats
}
