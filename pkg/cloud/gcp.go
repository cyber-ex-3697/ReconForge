package cloud

import (
    "fmt"
    "net/http"
)

type GCPDetector struct {
    timeout int
}

func NewGCPDetector(timeout int) *GCPDetector {
    return &GCPDetector{
        timeout: timeout,
    }
}

func (g *GCPDetector) CheckBucket(bucketName string) (*BucketResult, error) {
    url := fmt.Sprintf("https://storage.googleapis.com/%s", bucketName)
    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    return &BucketResult{
        URL:      url,
        Provider: "Google Cloud Storage",
        Public:   resp.StatusCode == 200,
    }, nil
}

func (g *GCPDetector) CheckFromDomain(domain string) ([]BucketResult, error) {
    patterns := []string{
        domain,
        domain + "-static",
        domain + "-assets",
        domain + "-cdn",
        "static." + domain,
        "assets." + domain,
        "cdn." + domain,
        "gcs." + domain,
    }
    
    var results []BucketResult
    for _, bucket := range patterns {
        result, err := g.CheckBucket(bucket)
        if err == nil {
            results = append(results, *result)
        }
    }
    return results, nil
}

func (g *GCPDetector) IsPublic(result *BucketResult) bool {
    return result.Public
}

func (g *GCPDetector) GetBucketACL(bucketName string) (string, error) {
    url := fmt.Sprintf("https://storage.googleapis.com/%s?acl", bucketName)
    resp, err := http.Get(url)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode == 200 {
        return "public", nil
    }
    return "private", nil
}
