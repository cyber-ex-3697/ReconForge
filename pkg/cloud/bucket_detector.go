package cloud

import (
    "fmt"
    "net/http"
    "strings"
)

type BucketResult struct {
    URL      string
    Provider string
    Public   bool
}

type BucketDetector struct {
    timeout int
}

func NewBucketDetector(timeout int) *BucketDetector {
    return &BucketDetector{
        timeout: timeout,
    }
}

func (b *BucketDetector) DetectAWS(domain string) ([]BucketResult, error) {
    var results []BucketResult
    
    patterns := []string{
        domain,
        domain + "-static",
        domain + "-assets",
        domain + "-cdn",
        "static." + domain,
        "assets." + domain,
        "cdn." + domain,
    }
    
    for _, bucket := range patterns {
        url := fmt.Sprintf("https://%s.s3.amazonaws.com", bucket)
        resp, err := http.Get(url)
        if err == nil {
            defer resp.Body.Close()
            results = append(results, BucketResult{
                URL:      url,
                Provider: "AWS S3",
                Public:   resp.StatusCode == 200,
            })
        }
    }
    return results, nil
}

func (b *BucketDetector) DetectGCP(domain string) ([]BucketResult, error) {
    var results []BucketResult
    
    patterns := []string{
        domain,
        domain + "-static",
        domain + "-assets",
    }
    
    for _, bucket := range patterns {
        url := fmt.Sprintf("https://storage.googleapis.com/%s", bucket)
        resp, err := http.Get(url)
        if err == nil {
            defer resp.Body.Close()
            results = append(results, BucketResult{
                URL:      url,
                Provider: "Google Cloud Storage",
                Public:   resp.StatusCode == 200,
            })
        }
    }
    return results, nil
}

func (b *BucketDetector) DetectAzure(domain string) ([]BucketResult, error) {
    var results []BucketResult
    
    patterns := []string{
        domain,
        domain + "static",
        domain + "assets",
    }
    
    for _, bucket := range patterns {
        url := fmt.Sprintf("https://%s.blob.core.windows.net", bucket)
        resp, err := http.Get(url)
        if err == nil {
            defer resp.Body.Close()
            results = append(results, BucketResult{
                URL:      url,
                Provider: "Azure Blob Storage",
                Public:   resp.StatusCode == 200,
            })
        }
    }
    return results, nil
}

func (b *BucketDetector) DetectAll(domain string) ([]BucketResult, error) {
    var allResults []BucketResult
    
    if results, err := b.DetectAWS(domain); err == nil {
        allResults = append(allResults, results...)
    }
    
    if results, err := b.DetectGCP(domain); err == nil {
        allResults = append(allResults, results...)
    }
    
    if results, err := b.DetectAzure(domain); err == nil {
        allResults = append(allResults, results...)
    }
    
    return allResults, nil
}
