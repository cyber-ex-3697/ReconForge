package benchmark

import (
    "runtime"
    "testing"
    
    "reconforge/internal/config"
    "reconforge/internal/logger"
)

func BenchmarkMemoryUsage(b *testing.B) {
    var memStats runtime.MemStats
    
    for i := 0; i < b.N; i++ {
        cfg := config.DefaultConfig()
        log := logger.New(false)
        
        runtime.GC()
        runtime.ReadMemStats(&memStats)
        
        _ = cfg
        _ = log
    }
}

func BenchmarkLargeConfig(b *testing.B) {
    for i := 0; i < b.N; i++ {
        cfg := config.DefaultConfig()
        cfg.Scan.Threads = 500
        cfg.Scan.DeepMode = true
        
        _ = cfg
    }
}
