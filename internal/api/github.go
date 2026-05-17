package api

import (
    "encoding/json"
    "fmt"
)

type GitHubResponse struct {
    Items []struct {
        HTMLURL string `json:"html_url"`
    } `json:"items"`
}

func (c *APIClient) GetGitHubSubdomains(domain string) ([]string, error) {
    token := c.GetAPIKey("github")
    if token == "" {
        return nil, fmt.Errorf("github token not configured")
    }
    
    url := fmt.Sprintf("https://api.github.com/search/code?q=%s&per_page=100", domain)
    headers := map[string]string{
        "Authorization": "token " + token,
        "Accept":        "application/vnd.github.v3+json",
    }
    
    data, err := c.doRequest(url, headers)
    if err != nil {
        return nil, err
    }
    
    var response GitHubResponse
    if err := json.Unmarshal(data, &response); err != nil {
        return nil, err
    }
    
    var urls []string
    for _, item := range response.Items {
        urls = append(urls, item.HTMLURL)
    }
    return urls, nil
}
