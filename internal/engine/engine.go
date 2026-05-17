package engine

import (
    "context"
    "fmt"
    "sync"
    "time"
    
    "reconforge/internal/config"
    "reconforge/internal/logger"
)

type Engine struct {
    config     *config.Config
    logger     *logger.Logger
    orchestrator *Orchestrator
    state      *ScanState
    mu         sync.RWMutex
}

type ScanResult struct {
    Target         string
    Subdomains     []string
    LiveHosts      []string
    URLs           []string
    Vulnerabilities []string
    StartTime      time.Time
    EndTime        time.Time
    Duration       time.Duration
}

func NewEngine(cfg *config.Config, log *logger.Logger) *Engine {
    return &Engine{
        config: cfg,
        logger: log,
        state:  NewScanState(),
    }
}

func (e *Engine) Run(ctx context.Context, target string) (*ScanResult, error) {
    e.logger.Info(fmt.Sprintf("Starting scan for target: %s", target))
    
    result := &ScanResult{
        Target:    target,
        StartTime: time.Now(),
    }
    
    // Initialize orchestrator
    e.orchestrator = NewOrchestrator(e.config, e.logger, e.state)
    
    // Phase 1: Subdomain Enumeration
    if err := e.orchestrator.RunPhase1(ctx, target, result); err != nil {
        e.logger.Error(fmt.Sprintf("Phase 1 failed: %v", err))
        return nil, err
    }
    
    // Phase 2: Live Host Detection
    if len(result.Subdomains) > 0 {
        if err := e.orchestrator.RunPhase2(ctx, result.Subdomains, result); err != nil {
            e.logger.Error(fmt.Sprintf("Phase 2 failed: %v", err))
        }
    }
    
    // Phase 3: URL Discovery
    if err := e.orchestrator.RunPhase3(ctx, target, result); err != nil {
        e.logger.Error(fmt.Sprintf("Phase 3 failed: %v", err))
    }
    
    // Phase 4: Vulnerability Scan (if deep mode)
    if e.config.Scan.DeepMode && len(result.LiveHosts) > 0 {
        if err := e.orchestrator.RunPhase4(ctx, result.LiveHosts, result); err != nil {
            e.logger.Error(fmt.Sprintf("Phase 4 failed: %v", err))
        }
    }
    
    result.EndTime = time.Now()
    result.Duration = result.EndTime.Sub(result.StartTime)
    
    e.logger.Success(fmt.Sprintf("Scan completed in %v", result.Duration))
    
    return result, nil
}

func (e *Engine) GetState() *ScanState {
    e.mu.RLock()
    defer e.mu.RUnlock()
    return e.state
}
