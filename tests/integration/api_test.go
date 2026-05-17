package integration

import (
    "testing"
    
    "reconforge/internal/api"
)

func TestNewAPIClient(t *testing.T) {
    client := api.NewAPIClient()
    
    if client == nil {
        t.Error("NewAPIClient returned nil")
    }
}

func TestAPIClientSetKey(t *testing.T) {
    client := api.NewAPIClient()
    client.SetAPIKey("test", "12345")
    
    key := client.GetAPIKey("test")
    if key != "12345" {
        t.Errorf("Expected key 12345, got %s", key)
    }
}

func TestEncryption(t *testing.T) {
    original := "test_secret_key"
    
    encrypted, err := api.Encrypt(original)
    if err != nil {
        t.Errorf("Encrypt failed: %v", err)
    }
    
    decrypted, err := api.Decrypt(encrypted)
    if err != nil {
        t.Errorf("Decrypt failed: %v", err)
    }
    
    if decrypted != original {
        t.Errorf("Expected %s, got %s", original, decrypted)
    }
}
