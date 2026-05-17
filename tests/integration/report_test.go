package integration

import (
    "os"
    "testing"
    
    "reconforge/pkg/report"
)

func TestNewGenerator(t *testing.T) {
    tmpDir := "/tmp/test_report"
    os.MkdirAll(tmpDir, 0755)
    defer os.RemoveAll(tmpDir)
    
    gen := report.NewGenerator(tmpDir)
    
    if gen == nil {
        t.Error("NewGenerator returned nil")
    }
}

func TestNewJSONExporter(t *testing.T) {
    exporter := report.NewJSONExporter("/tmp/test.json")
    
    if exporter == nil {
        t.Error("NewJSONExporter returned nil")
    }
}

func TestNewMarkdownExporter(t *testing.T) {
    exporter := report.NewMarkdownExporter("/tmp/test.md")
    
    if exporter == nil {
        t.Error("NewMarkdownExporter returned nil")
    }
}
