package cloud

import (
    "fmt"
    "net/http"
)

type AzureDetector struct {
    timeout int
}

func NewAzureDetector(timeout int) *AzureDetector {
    return &AzureDetector{
        timeout: timeout,
    }
}

func (a *AzureDetector) CheckContainer(containerName string) (*BucketResult, error) {
    url := fmt.Sprintf("https://%s.blob.core.windows.net", containerName)
    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    return &BucketResult{
        URL:      url,
        Provider: "Azure Blob Storage",
        Public:   resp.StatusCode == 200,
    }, nil
}

func (a *AzureDetector) CheckFromDomain(domain string) ([]BucketResult, error) {
    patterns := []string{
        domain,
        domain + "static",
        domain + "assets",
        domain + "cdn",
        "static" + domain,
        "assets" + domain,
        "cdn" + domain,
    }
    
    var results []BucketResult
    for _, container := range patterns {
        result, err := a.CheckContainer(container)
        if err == nil {
            results = append(results, *result)
        }
    }
    return results, nil
}

func (a *AzureDetector) IsPublic(result *BucketResult) bool {
    return result.Public
}

func (a *AzureDetector) CheckContainerACL(containerName string) (string, error) {
    url := fmt.Sprintf("https://%s.blob.core.windows.net?restype=container&comp=acl", containerName)
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
