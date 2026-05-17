package scheduler

import (
    "sync"
    "time"
)

type RetryableJob struct {
    Job       Job
    Attempts  int
    MaxRetries int
    LastError error
    NextRetry time.Time
}

type RetryQueue struct {
    queue    []*RetryableJob
    mu       sync.Mutex
    maxSize  int
}

func NewRetryQueue(maxSize int) *RetryQueue {
    return &RetryQueue{
        queue:   make([]*RetryableJob, 0),
        maxSize: maxSize,
    }
}

func (q *RetryQueue) Add(job Job, maxRetries int) {
    q.mu.Lock()
    defer q.mu.Unlock()
    
    if len(q.queue) >= q.maxSize {
        return
    }
    
    q.queue = append(q.queue, &RetryableJob{
        Job:        job,
        Attempts:   0,
        MaxRetries: maxRetries,
        NextRetry:  time.Now(),
    })
}

func (q *RetryQueue) MarkFailed(jobID int, err error) {
    q.mu.Lock()
    defer q.mu.Unlock()
    
    for _, rj := range q.queue {
        if rj.Job.ID == jobID {
            rj.Attempts++
            rj.LastError = err
            if rj.Attempts < rj.MaxRetries {
                backoff := time.Duration(1<<rj.Attempts) * time.Second
                rj.NextRetry = time.Now().Add(backoff)
            }
            break
        }
    }
}

func (q *RetryQueue) GetReadyJobs() []Job {
    q.mu.Lock()
    defer q.mu.Unlock()
    
    var ready []Job
    var remaining []*RetryableJob
    
    now := time.Now()
    for _, rj := range q.queue {
        if rj.Attempts >= rj.MaxRetries {
            continue
        }
        if rj.NextRetry.Before(now) {
            ready = append(ready, rj.Job)
        } else {
            remaining = append(remaining, rj)
        }
    }
    
    q.queue = remaining
    return ready
}

func (q *RetryQueue) Size() int {
    q.mu.Lock()
    defer q.mu.Unlock()
    return len(q.queue)
}

func (q *RetryQueue) Clear() {
    q.mu.Lock()
    defer q.mu.Unlock()
    q.queue = make([]*RetryableJob, 0)
}
