package unit

import (
    "os"
    "testing"
    
    "reconforge/internal/config"
)

func TestDefaultConfig(t *testing.T) {
    cfg := config.DefaultConfig()
    
    if cfg == nil {
        t.Error("DefaultConfig returned nil")
    }
    
    if cfg.Version != "4.0.0" {
        t.Errorf("Expected version 4.0.0, got %s", cfg.Version)
    }
    
    if cfg.Scan.Threads != 50 {
        t.Errorf("Expected threads 50, got %d", cfg.Scan.Threads)
    }
    
    if cfg.Scan.Timeout != 30 {
        t.Errorf("Expected timeout 30, got %d", cfg.Scan.Timeout)
    }
}

func TestLoadConfig(t *testing.T) {
    // Create temp config file
    tmpFile := "/tmp/test_config.yaml"
    content := `
version: "5.0.0"
scan:
  profile: "test"
  threads: 100
  deep_mode: true
  timeout: 60
  retries: 5
  rate_limit: 20
`
    os.WriteFile(tmpFile, []byte(content), 0644)
    defer os.Remove(tmpFile)
    
    cfg, err := config.Load(tmpFile)
    if err != nil {
        t.Errorf("Load failed: %v", err)
    }
    
    if cfg.Version != "5.0.0" {
        t.Errorf("Expected version 5.0.0, got %s", cfg.Version)
    }
    
    if cfg.Scan.Threads != 100 {
        t.Errorf("Expected threads 100, got %d", cfg.Scan.Threads)
    }
    
    if !cfg.Scan.DeepMode {
        t.Error("Expected deep_mode true")
    }
}

func TestConfigValidation(t *testing.T) {
    cfg := config.DefaultConfig()
    
    err := cfg.Validate()
    if err != nil {
        t.Errorf("Validation failed: %v", err)
    }
    
    // Test invalid threads
    cfg.Scan.Threads = 1000
    err = cfg.Validate()
    if err == nil {
        t.Error("Expected validation error for threads > 500")
    }
}
