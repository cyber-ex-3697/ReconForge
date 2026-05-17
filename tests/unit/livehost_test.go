package unit

import (
    "testing"
    
    "reconforge/pkg/livehost"
)

func TestNewProber(t *testing.T) {
    prober := livehost.NewProber(50)
    
    if prober == nil {
        t.Error("NewProber returned nil")
    }
}

func TestNewHTTPXWrapper(t *testing.T) {
    wrapper := livehost.NewHTTPXWrapper(50)
    
    if wrapper == nil {
        t.Error("NewHTTPXWrapper returned nil")
    }
}

func TestNewTechDetector(t *testing.T) {
    detector := livehost.NewTechDetector(50)
    
    if detector == nil {
        t.Error("NewTechDetector returned nil")
    }
}

func TestNewWAFDetector(t *testing.T) {
    detector := livehost.NewWAFDetector(10)
    
    if detector == nil {
        t.Error("NewWAFDetector returned nil")
    }
}

func TestNewFaviconHasher(t *testing.T) {
    hasher := livehost.NewFaviconHasher(10)
    
    if hasher == nil {
        t.Error("NewFaviconHasher returned nil")
    }
}
