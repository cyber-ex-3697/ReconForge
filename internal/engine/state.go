package engine

import (
    "sync"
    "time"
)

type Phase string

const (
    PhaseInit     Phase = "init"
    PhaseSubdomain Phase = "subdomain"
    PhaseLiveHost Phase = "live_host"
    PhaseURL      Phase = "url"
    PhaseVuln     Phase = "vulnerability"
    PhaseComplete Phase = "complete"
)

type ScanState struct {
    CurrentPhase   Phase
    Progress       int
    TotalSubdomains int
    TotalLiveHosts  int
    TotalURLs       int
    TotalVulns      int
    StartTime       time.Time
    LastUpdate      time.Time
    mu              sync.RWMutex
}

func NewScanState() *ScanState {
    return &ScanState{
        CurrentPhase: PhaseInit,
        Progress:     0,
        StartTime:    time.Now(),
        LastUpdate:   time.Now(),
    }
}

func (s *ScanState) SetPhase(phase Phase) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.CurrentPhase = phase
    s.LastUpdate = time.Now()
}

func (s *ScanState) GetPhase() Phase {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.CurrentPhase
}

func (s *ScanState) SetProgress(progress int) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.Progress = progress
    s.LastUpdate = time.Now()
}

func (s *ScanState) GetProgress() int {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.Progress
}

func (s *ScanState) GetDuration() time.Duration {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return time.Since(s.StartTime)
}
