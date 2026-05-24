package distributed

import (
    "sync"
    "time"
)

// HeartbeatMonitor monitors worker health
type HeartbeatMonitor struct {
    coordinator *Coordinator
    interval    time.Duration
    timeout     time.Duration
    stopChan    chan bool
    mu          sync.RWMutex
    running     bool
}

// NewHeartbeatMonitor creates a new heartbeat monitor
func NewHeartbeatMonitor(coordinator *Coordinator, interval, timeout time.Duration) *HeartbeatMonitor {
    return &HeartbeatMonitor{
        coordinator: coordinator,
        interval:    interval,
        timeout:     timeout,
        stopChan:    make(chan bool),
        running:     false,
    }
}

// Start starts the heartbeat monitor
func (h *HeartbeatMonitor) Start() {
    h.mu.Lock()
    if h.running {
        h.mu.Unlock()
        return
    }
    h.running = true
    h.mu.Unlock()
    
    go h.run()
}

// Stop stops the heartbeat monitor
func (h *HeartbeatMonitor) Stop() {
    h.mu.Lock()
    if !h.running {
        h.mu.Unlock()
        return
    }
    h.running = false
    h.mu.Unlock()
    
    h.stopChan <- true
}

// run is the main monitor loop
func (h *HeartbeatMonitor) run() {
    ticker := time.NewTicker(h.interval)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            h.checkWorkers()
        case <-h.stopChan:
            return
        }
    }
}

// checkWorkers checks all workers for heartbeats
func (h *HeartbeatMonitor) checkWorkers() {
    workers := h.coordinator.GetWorkerStats()
    now := time.Now()
    
    for _, worker := range workers {
        if now.Sub(worker.LastSeen) > h.timeout {
            // Worker is dead, unregister it
            h.coordinator.UnregisterWorker(worker.ID)
        }
    }
}

// SendHeartbeat sends a heartbeat for a worker
func (h *HeartbeatMonitor) SendHeartbeat(workerID string) {
    h.coordinator.UpdateWorkerHeartbeat(workerID)
}

// HeartbeatSender sends periodic heartbeats
type HeartbeatSender struct {
    workerID    string
    monitor     *HeartbeatMonitor
    interval    time.Duration
    stopChan    chan bool
    running     bool
    mu          sync.Mutex
}

// NewHeartbeatSender creates a new heartbeat sender
func NewHeartbeatSender(workerID string, monitor *HeartbeatMonitor, interval time.Duration) *HeartbeatSender {
    return &HeartbeatSender{
        workerID: workerID,
        monitor:  monitor,
        interval: interval,
        stopChan: make(chan bool),
        running:  false,
    }
}

// Start starts sending heartbeats
func (h *HeartbeatSender) Start() {
    h.mu.Lock()
    if h.running {
        h.mu.Unlock()
        return
    }
    h.running = true
    h.mu.Unlock()
    
    go h.run()
}

// Stop stops sending heartbeats
func (h *HeartbeatSender) Stop() {
    h.mu.Lock()
    if !h.running {
        h.mu.Unlock()
        return
    }
    h.running = false
    h.mu.Unlock()
    
    h.stopChan <- true
}

// run is the main sender loop
func (h *HeartbeatSender) run() {
    ticker := time.NewTicker(h.interval)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            h.monitor.SendHeartbeat(h.workerID)
        case <-h.stopChan:
            return
        }
    }
}
