package unit

import (
    "testing"
    
    "reconforge/pkg/subdomain"
)

func TestNewEnumerator(t *testing.T) {
    enum := subdomain.NewEnumerator("test.com", 50)
    
    if enum == nil {
        t.Error("NewEnumerator returned nil")
    }
}

func TestEnumeratorSetWordlist(t *testing.T) {
    enum := subdomain.NewEnumerator("test.com", 50)
    enum.SetWordlist("/tmp/wordlist.txt")
    
    // No panic means success
}

func TestResolver(t *testing.T) {
    resolver := subdomain.NewResolver()
    
    if resolver == nil {
        t.Error("NewResolver returned nil")
    }
}

func TestPassiveSource(t *testing.T) {
    source := subdomain.NewPassiveSource("test.com")
    
    if source == nil {
        t.Error("NewPassiveSource returned nil")
    }
}

func TestBruteForcer(t *testing.T) {
    brute := subdomain.NewBruteForcer("test.com", "/tmp/wordlist.txt")
    
    if brute == nil {
        t.Error("NewBruteForcer returned nil")
    }
}

func TestRecursiveEnumerator(t *testing.T) {
    rec := subdomain.NewRecursiveEnumerator("test.com", 3)
    
    if rec == nil {
        t.Error("NewRecursiveEnumerator returned nil")
    }
}
