package cloud

import (
    "fmt"
    "net/http"
    "strings"
)

type AWSDetector struct {
    timeout int
}

func NewAWSDetector(timeout int) *AWSDetector {
    return &AWSDetector{
        timeout: timeout,
    }
}

func (a *AWSDetector) CheckBucket(bucketName string) (*BucketResult, error) {
    url := fmt.Sprintf("https://%s.s3.amazonaws.com", bucketName)
    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    return &BucketResult{
        URL:      url,
        Provider: "AWS S3",
        Public:   resp.StatusCode == 200,
    }, nil
}

func (a *AWSDetector) CheckFromDomain(domain string) ([]BucketResult, error) {
    patterns := []string{
        domain,
        domain + "-static",
        domain + "-assets",
        domain + "-cdn",
        "static." + domain,
        "assets." + domain,
        "cdn." + domain,
        "s3." + domain,
    }
    
    var results []BucketResult
    for _, bucket := range patterns {
        result, err := a.CheckBucket(bucket)
        if err == nil {
            results = append(results, *result)
        }
    }
    return results, nil
}

func (a *AWSDetector) IsPublic(bucketResult *BucketResult) bool {
    return bucketResult.Public
}

func (a *AWSDetector) ListBuckets(domain string) ([]string, error) {
    // This would require AWS credentials
    // For now, return common patterns
    return []string{
        domain,
        domain + "-static",
        domain + "-assets",
        domain + "-backup",
        domain + "-logs",
    }, nil
}
