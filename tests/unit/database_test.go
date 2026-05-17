package unit

import (
    "os"
    "testing"
    
    "reconforge/internal/database"
)

func TestNewDatabase(t *testing.T) {
    tmpFile := "/tmp/test.db"
    defer os.Remove(tmpFile)
    
    db, err := database.New(tmpFile)
    if err != nil {
        t.Errorf("Failed to create database: %v", err)
    }
    defer db.Close()
    
    if db == nil {
        t.Error("Database is nil")
    }
}

func TestCreateScan(t *testing.T) {
    tmpFile := "/tmp/test.db"
    defer os.Remove(tmpFile)
    
    db, err := database.New(tmpFile)
    if err != nil {
        t.Fatalf("Failed to create database: %v", err)
    }
    defer db.Close()
    
    scanID, err := db.CreateScan("test.com", "init", "running")
    if err != nil {
        t.Errorf("CreateScan failed: %v", err)
    }
    
    if scanID == 0 {
        t.Error("Expected non-zero scan ID")
    }
}

func TestGetScan(t *testing.T) {
    tmpFile := "/tmp/test.db"
    defer os.Remove(tmpFile)
    
    db, err := database.New(tmpFile)
    if err != nil {
        t.Fatalf("Failed to create database: %v", err)
    }
    defer db.Close()
    
    scanID, _ := db.CreateScan("test.com", "init", "running")
    
    scan, err := db.GetScan(scanID)
    if err != nil {
        t.Errorf("GetScan failed: %v", err)
    }
    
    if scan.Target != "test.com" {
        t.Errorf("Expected target test.com, got %s", scan.Target)
    }
}

func TestAddSubdomain(t *testing.T) {
    tmpFile := "/tmp/test.db"
    defer os.Remove(tmpFile)
    
    db, err := database.New(tmpFile)
    if err != nil {
        t.Fatalf("Failed to create database: %v", err)
    }
    defer db.Close()
    
    scanID, _ := db.CreateScan("test.com", "init", "running")
    
    err = db.AddSubdomain(scanID, "sub.test.com", true)
    if err != nil {
        t.Errorf("AddSubdomain failed: %v", err)
    }
    
    subs, err := db.GetSubdomains(scanID)
    if err != nil {
        t.Errorf("GetSubdomains failed: %v", err)
    }
    
    if len(subs) != 1 {
        t.Errorf("Expected 1 subdomain, got %d", len(subs))
    }
}
