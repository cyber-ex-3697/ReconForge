package api

import (
    "encoding/json"
    "fmt"
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

// Chaos API - Get subdomains from ProjectDiscovery
func (c *APIClient) GetChaosSubdomains(domain string) ([]string, error) {
    key := c.GetAPIKey("chaos")
    if key == "" {
        return nil, fmt.Errorf("chaos API key not configured")
    }
    
    url := fmt.Sprintf("https://chaos.projectdiscovery.io/api/v1/subdomains?domain=%s", domain)
    req, _ := http.NewRequest("GET", url, nil)
    req.Header.Set("Authorization", key)
    
    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var result struct {
        Subdomains []string `json:"subdomains"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }
    
    return result.Subdomains, nil
}

// Shodan API - Search for hosts
func (c *APIClient) GetShodanHosts(domain string) ([]string, error) {
    key := c.GetAPIKey("shodan")
    if key == "" {
        return nil, fmt.Errorf("shodan API key not configured")
    }
    
    url := fmt.Sprintf("https://api.shodan.io/shodan/host/search?key=%s&query=hostname:%s", key, domain)
    resp, err := c.httpClient.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var result struct {
        Matches []struct {
            Hostnames []string `json:"hostnames"`
        } `json:"matches"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }
    
    var hosts []string
    for _, match := range result.Matches {
        hosts = append(hosts, match.Hostnames...)
    }
    return hosts, nil
}

// GitHub API - Search code for subdomains
func (c *APIClient) GetGitHubSubdomains(domain string) ([]string, error) {
    token := c.GetAPIKey("github")
    if token == "" {
        return nil, fmt.Errorf("github token not configured")
    }
    
    url := fmt.Sprintf("https://api.github.com/search/code?q=%s&per_page=100", domain)
    req, _ := http.NewRequest("GET", url, nil)
    req.Header.Set("Authorization", "token "+token)
    
    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var result struct {
        Items []struct {
            HTMLURL string `json:"html_url"`
        } `json:"items"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }
    
    var urls []string
    for _, item := range result.Items {
        urls = append(urls, item.HTMLURL)
    }
    return urls, nil
}
