package distributed

import (
    "sync"
    "time"
)

type Worker struct {
    id       string
    maxConcurrent int
    semaphore chan struct{}
}

func NewWorker(id string, maxConcurrent int) *Worker {
    return &Worker{
        id:           id,
        maxConcurrent: maxConcurrent,
        semaphore:    make(chan struct{}, maxConcurrent),
    }
}

func (w *Worker) Run(task func()) {
    w.semaphore <- struct{}{}
    go func() {
        defer func() { <-w.semaphore }()
        task()
    }()
}

func (w *Worker) Wait() {
    for i := 0; i < cap(w.semaphore); i++ {
        w.semaphore <- struct{}{}
    }
}
