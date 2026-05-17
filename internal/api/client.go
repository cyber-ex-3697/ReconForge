package api

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "time"
)

type APIClient struct {
    httpClient *http.Client
    apiKeys    map[string]string
}

func NewAPIClient() *APIClient {
    return &APIClient{
        httpClient: &http.Client{Timeout: 30 * time.Second},
        apiKeys:    make(map[string]string),
    }
}

func (c *APIClient) SetAPIKey(service, key string) {
    c.apiKeys[service] = key
}

func (c *APIClient) GetAPIKey(service string) string {
    return c.apiKeys[service]
}

func (c *APIClient) doRequest(url string, headers map[string]string) ([]byte, error) {
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, err
    }
    
    for k, v := range headers {
        req.Header.Set(k, v)
    }
    
    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != 200 {
        return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
    }
    
    return io.ReadAll(resp.Body)
}
