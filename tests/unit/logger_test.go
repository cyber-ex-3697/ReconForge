package unit

import (
    "bytes"
    "testing"
    
    "reconforge/internal/logger"
)

func TestNewLogger(t *testing.T) {
    log := logger.New(true)
    
    if log == nil {
        t.Error("NewLogger returned nil")
    }
}

func TestLoggerInfo(t *testing.T) {
    var buf bytes.Buffer
    // Note: Would need to capture output for full test
    
    log := logger.New(true)
    
    // Just verify no panic
    defer func() {
        if r := recover(); r != nil {
            t.Errorf("Info panicked: %v", r)
        }
    }()
    
    log.Info("Test message")
    log.Debug("Debug message")
    log.Warn("Warning message")
    log.Error("Error message")
    log.Success("Success message")
}

func TestLoggerLevels(t *testing.T) {
    log := logger.New(false) // debug disabled
    
    // Should not panic
    log.Debug("This should not be printed")
    log.Info("This should be printed")
}
