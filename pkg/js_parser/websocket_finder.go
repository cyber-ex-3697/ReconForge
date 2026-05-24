package js_parser

import (
    "regexp"
    "strings"
)

// WebSocketInfo represents WebSocket connection info
type WebSocketInfo struct {
    URL      string            `json:"url"`
    Protocols []string         `json:"protocols"`
    Events   map[string]string `json:"events"`
    Raw      string            `json:"raw"`
}

// WebsocketFinder finds WebSocket connections in JavaScript
type WebsocketFinder struct {
    websockets []WebSocketInfo
    patterns   []*regexp.Regexp
}

// NewWebsocketFinder creates a new WebSocket finder
func NewWebsocketFinder() *WebsocketFinder {
    return &WebsocketFinder{
        websockets: make([]WebSocketInfo, 0),
        patterns: []*regexp.Regexp{
            regexp.MustCompile(`new\s+WebSocket\s*\(\s*['"]([^'"]+)['"]`),
            regexp.MustCompile(`WebSocket\s*\.\s*send\s*\(\s*['"]([^'"]+)['"]`),
            regexp.MustCompile(`ws://[a-zA-Z0-9./_-]+`),
            regexp.MustCompile(`wss://[a-zA-Z0-9./_-]+`),
            regexp.MustCompile(`socket\.io\s*\(\s*['"]([^'"]+)['"]`),
            regexp.MustCompile(`io\s*\(\s*['"]([^'"]+)['"]`),
            regexp.MustCompile(`['"]wss?://[^'"]+['"]`),
        },
    }
}

// Find searches for WebSocket connections in source
func (f *WebsocketFinder) Find(source string) []WebSocketInfo {
    f.websockets = make([]WebSocketInfo, 0)
    
    for _, pattern := range f.patterns {
        matches := pattern.FindAllStringSubmatch(source, -1)
        for _, match := range matches {
            url := match[0]
            if len(match) > 1 {
                url = match[1]
            }
            
            info := WebSocketInfo{
                URL:      f.cleanURL(url),
                Protocols: f.extractProtocols(source, url),
                Events:   f.extractEvents(source, url),
                Raw:      match[0],
            }
            f.websockets = append(f.websockets, info)
        }
    }
    
    // Deduplicate
    return f.deduplicate()
}

// cleanURL cleans up WebSocket URL
func (f *WebsocketFinder) cleanURL(url string) string {
    url = strings.Trim(url, "\"'` ")
    
    // Ensure protocol
    if !strings.HasPrefix(url, "ws://") && !strings.HasPrefix(url, "wss://") {
        if strings.HasPrefix(url, "//") {
            url = "wss:" + url
        } else if strings.HasPrefix(url, "/") {
            // Relative path - would need base URL
        }
    }
    
    return url
}

// extractProtocols extracts protocols from WebSocket constructor
func (f *WebsocketFinder) extractProtocols(source, url string) []string {
    protocols := make([]string, 0)
    
    // Look for second argument in WebSocket constructor
    pattern := regexp.MustCompile(`new\s+WebSocket\s*\(\s*['"][^'"]+['"]\s*,\s*\[([^\]]+)\]`)
    matches := pattern.FindStringSubmatch(source)
    if len(matches) > 1 {
        for _, p := range strings.Split(matches[1], ",") {
            p = strings.Trim(p, ` "'`)
            if p != "" {
                protocols = append(protocols, p)
            }
        }
    }
    
    return protocols
}

// extractEvents extracts event handlers for WebSocket
func (f *WebsocketFinder) extractEvents(source, url string) map[string]string {
    events := make(map[string]string)
    
    eventTypes := []string{"onopen", "onmessage", "onerror", "onclose"}
    for _, event := range eventTypes {
        pattern := regexp.MustCompile(`\.` + event + `\s*=\s*function\s*\([^)]*\)\s*{([^}]+)}`)
        matches := pattern.FindStringSubmatch(source)
        if len(matches) > 1 {
            events[event] = strings.TrimSpace(matches[1])
        }
    }
    
    return events
}

// deduplicate removes duplicate WebSocket connections
func (f *WebsocketFinder) deduplicate() []WebSocketInfo {
    seen := make(map[string]bool)
    var unique []WebSocketInfo
    
    for _, ws := range f.websockets {
        if !seen[ws.URL] {
            seen[ws.URL] = true
            unique = append(unique, ws)
        }
    }
    
    return unique
}

// GetWebsockets returns all found WebSocket connections
func (f *WebsocketFinder) GetWebsockets() []WebSocketInfo {
    return f.websockets
}
