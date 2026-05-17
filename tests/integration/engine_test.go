package integration

import (
    "context"
    "testing"
    
    "reconforge/internal/config"
    "reconforge/internal/engine"
    "reconforge/internal/logger"
)

func TestNewEngine(t *testing.T) {
    cfg := config.DefaultConfig()
    log := logger.New(true)
    
    eng := engine.NewEngine(cfg, log)
    
    if eng == nil {
        t.Error("NewEngine returned nil")
    }
}

func TestEngineGetState(t *testing.T) {
    cfg := config.DefaultConfig()
    log := logger.New(true)
    
    eng := engine.NewEngine(cfg, log)
    state := eng.GetState()
    
    if state == nil {
        t.Error("GetState returned nil")
    }
}
