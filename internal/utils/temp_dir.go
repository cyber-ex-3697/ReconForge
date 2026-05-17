package utils

import (
    "fmt"
    "os"
    "path/filepath"
)

func CreateTempDir() (string, error) {
    tempDir := filepath.Join(os.TempDir(), fmt.Sprintf("reconforge_%d", os.Getpid()))
    if err := os.MkdirAll(tempDir, 0700); err != nil {
        return "", err
    }
    return tempDir, nil
}

func CleanupTempDir(tempDir string) error {
    if tempDir != "" {
        return os.RemoveAll(tempDir)
    }
    return nil
}

func CreateSecureTempFile(prefix string) (*os.File, error) {
    return os.CreateTemp("", prefix+"_*.tmp")
}
