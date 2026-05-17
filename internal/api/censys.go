package api

import (
    "encoding/json"
    "fmt"
)

type CensysResponse struct {
    Results []struct {
        Name string `json:"name"`
    } `json:"results"`
}

func (c *APIClient) GetCensysSubdomains(domain string) ([]string, error) {
    apiID := c.GetAPIKey("censys_id")
    apiSecret := c.GetAPIKey("censys_secret")
    
    if apiID == "" || apiSecret == "" {
        return nil, fmt.Errorf("censys API credentials not configured")
    }
    
    url := fmt.Sprintf("https://search.censys.io/api/v2/certificates/search?q=%s&per_page=100", domain)
    
    data, err := c.doRequest(url, nil)
    if err != nil {
        return nil, err
    }
    
    var response CensysResponse
    if err := json.Unmarshal(data, &response); err != nil {
        return nil, err
    }
    
    var subdomains []string
    for _, result := range response.Results {
        subdomains = append(subdomains, result.Name)
    }
    return subdomains, nil
}
