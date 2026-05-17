package benchmark

import (
    "testing"
    
    "reconforge/pkg/subdomain"
)

func BenchmarkSubdomainEnumeration(b *testing.B) {
    enum := subdomain.NewEnumerator("example.com", 50)
    
    for i := 0; i < b.N; i++ {
        _, _ = enum.Run()
    }
}

func BenchmarkResolver(b *testing.B) {
    resolver := subdomain.NewResolver()
    domains := []string{"example.com", "google.com", "github.com"}
    
    for i := 0; i < b.N; i++ {
        _, _ = resolver.Resolve(domains)
    }
}
