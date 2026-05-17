package benchmark

import (
    "testing"
    
    "reconforge/internal/config"
)

func BenchmarkConfigLoad(b *testing.B) {
    for i := 0; i < b.N; i++ {
        _ = config.DefaultConfig()
    }
}

func BenchmarkConfigValidation(b *testing.B) {
    cfg := config.DefaultConfig()
    
    for i := 0; i < b.N; i++ {
        _ = cfg.Validate()
    }
}
