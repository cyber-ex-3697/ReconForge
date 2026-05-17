package api

import (
    "encoding/json"
    "fmt"
)

type SecurityTrailsResponse struct {
    Subdomains []string `json:"subdomains"`
}

func (c *APIClient) GetSecurityTrailsSubdomains(domain string) ([]string, error) {
    key := c.GetAPIKey("securitytrails")
    if key == "" {
        return nil, fmt.Errorf("securitytrails API key not configured")
    }
    
    url := fmt.Sprintf("https://api.securitytrails.com/v1/domain/%s/subdomains", domain)
    headers := map[string]string{
        "APIKEY": key,
    }
    
    data, err := c.doRequest(url, headers)
    if err != nil {
        return nil, err
    }
    
    var response SecurityTrailsResponse
    if err := json.Unmarshal(data, &response); err != nil {
        return nil, err
    }
    
    var result []string
    for _, sub := range response.Subdomains {
        result = append(result, sub+"."+domain)
    }
    return result, nil
}
