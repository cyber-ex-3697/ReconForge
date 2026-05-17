package utils

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "io"
    "net/http"
    "os"
)

func CalculateSHA256(data []byte) string {
    hash := sha256.Sum256(data)
    return hex.EncodeToString(hash[:])
}

func CalculateFileSHA256(filepath string) (string, error) {
    file, err := os.Open(filepath)
    if err != nil {
        return "", err
    }
    defer file.Close()
    
    hash := sha256.New()
    if _, err := io.Copy(hash, file); err != nil {
        return "", err
    }
    
    return hex.EncodeToString(hash.Sum(nil)), nil
}

func VerifyChecksum(filepath, expectedChecksum string) (bool, error) {
    actual, err := CalculateFileSHA256(filepath)
    if err != nil {
        return false, err
    }
    return actual == expectedChecksum, nil
}

func DownloadWithChecksum(url, destPath, expectedChecksum string) error {
    resp, err := http.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    file, err := os.Create(destPath)
    if err != nil {
        return err
    }
    defer file.Close()
    
    hash := sha256.New()
    writer := io.MultiWriter(file, hash)
    
    if _, err := io.Copy(writer, resp.Body); err != nil {
        return err
    }
    
    if expectedChecksum != "" {
        actualChecksum := hex.EncodeToString(hash.Sum(nil))
        if actualChecksum != expectedChecksum {
            return fmt.Errorf("checksum mismatch: expected %s, got %s", expectedChecksum, actualChecksum)
        }
    }
    
    return nil
}
