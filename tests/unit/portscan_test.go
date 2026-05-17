package unit

import (
    "testing"
    
    "reconforge/pkg/portscan"
)

func TestNewScanner(t *testing.T) {
    scanner := portscan.NewScanner(50)
    
    if scanner == nil {
        t.Error("NewScanner returned nil")
    }
}

func TestScannerSetPorts(t *testing.T) {
    scanner := portscan.NewScanner(50)
    scanner.SetPorts("80,443,22")
    
    // No panic means success
}

func TestNewNaabuWrapper(t *testing.T) {
    wrapper := portscan.NewNaabuWrapper(1000, 50)
    
    if wrapper == nil {
        t.Error("NewNaabuWrapper returned nil")
    }
}

func TestNewNmapWrapper(t *testing.T) {
    wrapper := portscan.NewNmapWrapper()
    
    if wrapper == nil {
        t.Error("NewNmapWrapper returned nil")
    }
}

func TestNewMasscanWrapper(t *testing.T) {
    wrapper := portscan.NewMasscanWrapper(1000)
    
    if wrapper == nil {
        t.Error("NewMasscanWrapper returned nil")
    }
}
