package livehost

import (
    "crypto/md5"
    "encoding/hex"
    "fmt"
    "io"
    "net/http"
    "strings"
)

type FaviconHasher struct {
    timeout int
}

func NewFaviconHasher(timeout int) *FaviconHasher {
    return &FaviconHasher{
        timeout: timeout,
    }
}

func (f *FaviconHasher) GetFaviconHash(url string) (string, error) {
    paths := []string{
        "/favicon.ico",
        "/favicon.png",
        "/favicon.jpg",
        "/static/favicon.ico",
        "/assets/favicon.ico",
        "/images/favicon.ico",
    }
    
    for _, path := range paths {
        fullURL := strings.TrimSuffix(url, "/") + path
        resp, err := http.Get(fullURL)
        if err != nil {
            continue
        }
        defer resp.Body.Close()
        
        if resp.StatusCode == 200 {
            body, err := io.ReadAll(resp.Body)
            if err != nil {
                continue
            }
            hash := md5.Sum(body)
            return hex.EncodeToString(hash[:]), nil
        }
    }
    return "", fmt.Errorf("no favicon found for %s", url)
}

func (f *FaviconHasher) GetFaviconHashBatch(urls []string) (map[string]string, error) {
    results := make(map[string]string)
    for _, url := range urls {
        if hash, err := f.GetFaviconHash(url); err == nil {
            results[url] = hash
        }
    }
    return results, nil
}

func (f *FaviconHasher) CompareHash(hash1, hash2 string) bool {
    return hash1 == hash2
}
