package scheduler

import (
    "context"
    "sync"
)

type Job struct {
    ID      int
    Task    func() error
    Retry   int
}

type Result struct {
    JobID int
    Err   error
}

type WorkerPool struct {
    workers     int
    jobQueue    chan Job
    resultQueue chan Result
    wg          sync.WaitGroup
    ctx         context.Context
    cancel      context.CancelFunc
}

func NewWorkerPool(workers int) *WorkerPool {
    ctx, cancel := context.WithCancel(context.Background())
    return &WorkerPool{
        workers:     workers,
        jobQueue:    make(chan Job, 10000),
        resultQueue: make(chan Result, 10000),
        ctx:         ctx,
        cancel:      cancel,
    }
}

func (p *WorkerPool) Start() {
    for i := 0; i < p.workers; i++ {
        p.wg.Add(1)
        go p.worker()
    }
}

func (p *WorkerPool) worker() {
    defer p.wg.Done()
    
    for {
        select {
        case <-p.ctx.Done():
            return
        case job, ok := <-p.jobQueue:
            if !ok {
                return
            }
            
            var err error
            for attempt := 0; attempt <= job.Retry; attempt++ {
                if attempt > 0 {
                    // Exponential backoff
                    // time.Sleep(time.Duration(1<<attempt) * time.Second)
                }
                err = job.Task()
                if err == nil {
                    break
                }
            }
            
            p.resultQueue <- Result{
                JobID: job.ID,
                Err:   err,
            }
        }
    }
}

func (p *WorkerPool) Submit(job Job) {
    select {
    case p.jobQueue <- job:
    case <-p.ctx.Done():
    }
}

func (p *WorkerPool) Results() <-chan Result {
    return p.resultQueue
}

func (p *WorkerPool) Stop() {
    p.cancel()
    close(p.jobQueue)
    p.wg.Wait()
    close(p.resultQueue)
}

func (p *WorkerPool) GetWorkerCount() int {
    return p.workers
}
