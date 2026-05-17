package logger

import (
    "os"
    "path/filepath"
    "sync"
)

type FileWriter struct {
    file   *os.File
    mu     sync.Mutex
    path   string
}

func NewFileWriter(path string) (*FileWriter, error) {
    dir := filepath.Dir(path)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return nil, err
    }
    
    file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
        return nil, err
    }
    
    return &FileWriter{
        file: file,
        path: path,
    }, nil
}

func (w *FileWriter) Write(p []byte) (int, error) {
    w.mu.Lock()
    defer w.mu.Unlock()
    return w.file.Write(p)
}

func (w *FileWriter) Close() error {
    return w.file.Close()
}
