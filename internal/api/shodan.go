package api

import (
    "encoding/json"
    "fmt"
)

type ShodanResponse struct {
    Matches []struct {
        Hostnames []string `json:"hostnames"`
        Ports     []int    `json:"ports"`
    } `json:"matches"`
}

func (c *APIClient) GetShodanHosts(domain string) ([]string, error) {
    key := c.GetAPIKey("shodan")
    if key == "" {
        return nil, fmt.Errorf("shodan API key not configured")
    }
    
    url := fmt.Sprintf("https://api.shodan.io/shodan/host/search?key=%s&query=hostname:%s", key, domain)
    
    data, err := c.doRequest(url, nil)
    if err != nil {
        return nil, err
    }
    
    var response ShodanResponse
    if err := json.Unmarshal(data, &response); err != nil {
        return nil, err
    }
    
    var hosts []string
    for _, match := range response.Matches {
        hosts = append(hosts, match.Hostnames...)
    }
    return hosts, nil
}
