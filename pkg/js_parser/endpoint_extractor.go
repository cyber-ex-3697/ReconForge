package js_parser

import (
    "regexp"
    "strings"
)

// EndpointExtractor extracts API endpoints from JavaScript
type EndpointExtractor struct {
    patterns    []*regexp.Regexp
    endpoints   map[string]bool
    baseURL     string
}

// NewEndpointExtractor creates a new endpoint extractor
func NewEndpointExtractor() *EndpointExtractor {
    return &EndpointExtractor{
        patterns: []*regexp.Regexp{
            regexp.MustCompile(`/api/[a-zA-Z0-9/_-]+`),
            regexp.MustCompile(`/v\d/[a-zA-Z0-9/_-]+`),
            regexp.MustCompile(`/rest/[a-zA-Z0-9/_-]+`),
            regexp.MustCompile(`/graphql/?`),
            regexp.MustCompile(`/auth/[a-zA-Z0-9/_-]+`),
            regexp.MustCompile(`/user/[a-zA-Z0-9/_-]+`),
            regexp.MustCompile(`/admin/[a-zA-Z0-9/_-]+`),
            regexp.MustCompile(`/upload/[a-zA-Z0-9/_-]+`),
            regexp.MustCompile(`/download/[a-zA-Z0-9/_-]+`),
            regexp.MustCompile(`/payment/[a-zA-Z0-9/_-]+`),
            regexp.MustCompile(`/[a-z]+/[0-9]+/[a-z]+`),
            regexp.MustCompile(`\?[a-zA-Z0-9_]+=[^&]+`),
        },
        endpoints: make(map[string]bool),
    }
}

// SetBaseURL sets the base URL for absolute endpoint generation
func (e *EndpointExtractor) SetBaseURL(baseURL string) {
    e.baseURL = strings.TrimSuffix(baseURL, "/")
}

// Extract extracts endpoints from JavaScript source
func (e *EndpointExtractor) Extract(source string) []string {
    e.endpoints = make(map[string]bool)
    
    for _, pattern := range e.patterns {
        matches := pattern.FindAllString(source, -1)
        for _, match := range matches {
            e.addEndpoint(match)
        }
    }
    
    // Extract from strings that look like routes
    routePattern := regexp.MustCompile(`['"](/[a-zA-Z0-9/_{}:?-]+)['"]`)
    matches := routePattern.FindAllStringSubmatch(source, -1)
    for _, match := range matches {
        if len(match) > 1 {
            e.addEndpoint(match[1])
        }
    }
    
    // Extract from fetch/axios calls
    fetchPattern := regexp.MustCompile(`(?:fetch|axios|request)\s*\(\s*['"]([^'"]+)['"]`)
    matches = fetchPattern.FindAllStringSubmatch(source, -1)
    for _, match := range matches {
        if len(match) > 1 {
            e.addEndpoint(match[1])
        }
    }
    
    result := make([]string, 0, len(e.endpoints))
    for ep := range e.endpoints {
        result = append(result, ep)
    }
    return result
}

// addEndpoint adds an endpoint to the collection
func (e *EndpointExtractor) addEndpoint(endpoint string) {
    if endpoint == "" {
        return
    }
    
    // Clean up the endpoint
    endpoint = strings.TrimSpace(endpoint)
    endpoint = strings.Trim(endpoint, "\"'`")
    
    // Convert to absolute URL if baseURL is set
    if e.baseURL != "" && strings.HasPrefix(endpoint, "/") {
        endpoint = e.baseURL + endpoint
    }
    
    // Remove query parameters for deduplication
    if idx := strings.Index(endpoint, "?"); idx != -1 {
        base := endpoint[:idx]
        params := endpoint[idx+1:]
        
        // Keep the base endpoint
        e.endpoints[base] = true
        
        // Also extract individual parameters
        for _, param := range strings.Split(params, "&") {
            if strings.Contains(param, "=") {
                parts := strings.SplitN(param, "=", 2)
                e.endpoints[base+"?"+parts[0]+"="] = true
            }
        }
    } else {
        e.endpoints[endpoint] = true
    }
}

// GetEndpoints returns all extracted endpoints
func (e *EndpointExtractor) GetEndpoints() []string {
    result := make([]string, 0, len(e.endpoints))
    for ep := range e.endpoints {
        result = append(result, ep)
    }
    return result
}

// GetUniqueAPIPaths returns unique API paths
func (e *EndpointExtractor) GetUniqueAPIPaths() []string {
    apiPaths := make(map[string]bool)
    
    for ep := range e.endpoints {
        if strings.Contains(ep, "/api/") || 
           strings.Contains(ep, "/v1/") || 
           strings.Contains(ep, "/v2/") ||
           strings.Contains(ep, "/graphql") {
            apiPaths[ep] = true
        }
    }
    
    result := make([]string, 0, len(apiPaths))
    for path := range apiPaths {
        result = append(result, path)
    }
    return result
}
