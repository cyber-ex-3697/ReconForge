package unit

import (
    "testing"
    
    "reconforge/pkg/url"
)

func TestNewCollector(t *testing.T) {
    collector := url.NewCollector("test.com")
    
    if collector == nil {
        t.Error("NewCollector returned nil")
    }
}

func TestNewGAUWrapper(t *testing.T) {
    gau := url.NewGAUWrapper()
    
    if gau == nil {
        t.Error("NewGAUWrapper returned nil")
    }
}

func TestNewKatanaWrapper(t *testing.T) {
    katana := url.NewKatanaWrapper(3, 50)
    
    if katana == nil {
        t.Error("NewKatanaWrapper returned nil")
    }
}

func TestNewJSExtractor(t *testing.T) {
    extractor := url.NewJSExtractor()
    
    if extractor == nil {
        t.Error("NewJSExtractor returned nil")
    }
}

func TestNewParamExtractor(t *testing.T) {
    extractor := url.NewParamExtractor()
    
    if extractor == nil {
        t.Error("NewParamExtractor returned nil")
    }
}

func TestNewAPIDetector(t *testing.T) {
    detector := url.NewAPIDetector()
    
    if detector == nil {
        t.Error("NewAPIDetector returned nil")
    }
}
