package unit

import (
    "testing"
    
    "reconforge/pkg/takeover"
)

func TestNewDetector(t *testing.T) {
    detector := takeover.NewDetector(20)
    
    if detector == nil {
        t.Error("NewDetector returned nil")
    }
}

func TestNewSubzyWrapper(t *testing.T) {
    wrapper := takeover.NewSubzyWrapper(20, 10)
    
    if wrapper == nil {
        t.Error("NewSubzyWrapper returned nil")
    }
}

func TestNewFingerprintMatcher(t *testing.T) {
    matcher := takeover.NewFingerprintMatcher()
    
    if matcher == nil {
        t.Error("NewFingerprintMatcher returned nil")
    }
}

func TestNewCNAMEChecker(t *testing.T) {
    checker := takeover.NewCNAMEChecker()
    
    if checker == nil {
        t.Error("NewCNAMEChecker returned nil")
    }
}

func TestFingerprintMatch(t *testing.T) {
    matcher := takeover.NewFingerprintMatcher()
    
    // Test matching
    match := matcher.Match("NoSuchBucket error message")
    if match == nil {
        t.Log("No match found (expected for local test)")
    }
}

func TestCNAMECheckerCheck(t *testing.T) {
    checker := takeover.NewCNAMEChecker()
    
    // Test with invalid domain (won't resolve)
    vulnerable, cname, service := checker.Check("nonexistent.example.com")
    
    // Should not panic
    _ = vulnerable
    _ = cname
    _ = service
}
