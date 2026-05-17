package url

import (
    "regexp"
    "strings"
)

type APIDetector struct {
    patterns []string
}

func NewAPIDetector() *APIDetector {
    return &APIDetector{
        patterns: []string{
            `api`,
            `v1`,
            `v2`,
            `v3`,
            `graphql`,
            `rest`,
            `swagger`,
            `openapi`,
            `wp-json`,
            `rpc`,
            `jsonrpc`,
            `/api/`,
            `/v1/`,
            `/v2/`,
            `/v3/`,
        },
    }
}

func (d *APIDetector) Detect(urls []string) []string {
    var apiEndpoints []string
    seen := make(map[string]bool)
    
    for _, url := range urls {
        lowerURL := strings.ToLower(url)
        for _, pattern := range d.patterns {
            if strings.Contains(lowerURL, pattern) {
                if !seen[url] {
                    seen[url] = true
                    apiEndpoints = append(apiEndpoints, url)
                }
                break
            }
        }
    }
    return apiEndpoints
}

func (d *APIDetector) DetectWithRegex(urls []string) []string {
    var apiEndpoints []string
    seen := make(map[string]bool)
    
    for _, pattern := range d.patterns {
        re := regexp.MustCompile(pattern)
        for _, url := range urls {
            if re.MatchString(url) {
                if !seen[url] {
                    seen[url] = true
                    apiEndpoints = append(apiEndpoints, url)
                }
            }
        }
    }
    return apiEndpoints
}

func (d *APIDetector) AddPattern(pattern string) {
    d.patterns = append(d.patterns, pattern)
}
