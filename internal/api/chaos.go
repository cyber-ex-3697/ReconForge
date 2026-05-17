package api

import (
    "encoding/json"
    "fmt"
)

type ChaosResponse struct {
    Subdomains []string `json:"subdomains"`
    Count      int      `json:"count"`
}

func (c *APIClient) GetChaosSubdomains(domain string) ([]string, error) {
    key := c.GetAPIKey("chaos")
    if key == "" {
        return nil, fmt.Errorf("chaos API key not configured")
    }
    
    url := fmt.Sprintf("https://chaos.projectdiscovery.io/api/v1/subdomains?domain=%s", domain)
    headers := map[string]string{
        "Authorization": key,
    }
    
    data, err := c.doRequest(url, headers)
    if err != nil {
        return nil, err
    }
    
    var response ChaosResponse
    if err := json.Unmarshal(data, &response); err != nil {
        return nil, err
    }
    
    return response.Subdomains, nil
}
