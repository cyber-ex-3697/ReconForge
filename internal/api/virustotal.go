package api

import (
    "encoding/json"
    "fmt"
)

type VirusTotalResponse struct {
    Data struct {
        Attributes struct {
            LastAnalysisStats struct {
                Malicious int `json:"malicious"`
            } `json:"last_analysis_stats"`
        } `json:"attributes"`
    } `json:"data"`
}

func (c *APIClient) CheckVirusTotal(domain string) (bool, error) {
    key := c.GetAPIKey("virustotal")
    if key == "" {
        return false, fmt.Errorf("virustotal API key not configured")
    }
    
    url := fmt.Sprintf("https://www.virustotal.com/api/v3/domains/%s", domain)
    headers := map[string]string{
        "x-apikey": key,
    }
    
    data, err := c.doRequest(url, headers)
    if err != nil {
        return false, err
    }
    
    var response VirusTotalResponse
    if err := json.Unmarshal(data, &response); err != nil {
        return false, err
    }
    
    return response.Data.Attributes.LastAnalysisStats.Malicious > 0, nil
}
