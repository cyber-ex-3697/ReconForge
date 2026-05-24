package js_parser

import (
    "regexp"
    "strings"
)

// Parameter represents a mined parameter
type Parameter struct {
    Name      string   `json:"name"`
    Type      string   `json:"type"` // query, path, body
    Examples  []string `json:"examples"`
    Required  bool     `json:"required"`
}

// ParamMiner mines parameters from JavaScript
type ParamMiner struct {
    parameters map[string]*Parameter
    patterns   []*regexp.Regexp
}

// NewParamMiner creates a new parameter miner
func NewParamMiner() *ParamMiner {
    return &ParamMiner{
        parameters: make(map[string]*Parameter),
        patterns: []*regexp.Regexp{
            regexp.MustCompile(`req\.params\.([a-zA-Z0-9_]+)`),
            regexp.MustCompile(`req\.query\.([a-zA-Z0-9_]+)`),
            regexp.MustCompile(`req\.body\.([a-zA-Z0-9_]+)`),
            regexp.MustCompile(`\?([a-zA-Z0-9_]+)=`),
            regexp.MustCompile(`&([a-zA-Z0-9_]+)=`),
            regexp.MustCompile(`params:\s*{([^}]+)}`),
            regexp.MustCompile(`data:\s*{([^}]+)}`),
            regexp.MustCompile(`query:\s*{([^}]+)}`),
            regexp.MustCompile(`:([a-zA-Z0-9_]+)\s*=>`),
            regexp.MustCompile(`get\('([^']+)'\)`),
            regexp.MustCompile(`post\('([^']+)'\)`),
        },
    }
}

// Mine extracts parameters from JavaScript source
func (p *ParamMiner) Mine(source string) []*Parameter {
    p.parameters = make(map[string]*Parameter)
    
    for _, pattern := range p.patterns {
        matches := pattern.FindAllStringSubmatch(source, -1)
        for _, match := range matches {
            if len(match) > 1 {
                p.addParameter(match[1], p.getParamType(pattern.String()))
            }
        }
    }
    
    // Extract from URL patterns
    urlPattern := regexp.MustCompile(`['"]/([a-zA-Z0-9/_:]+)['"]`)
    matches := urlPattern.FindAllStringSubmatch(source, -1)
    for _, match := range matches {
        if len(match) > 1 {
            parts := strings.Split(match[1], "/")
            for _, part := range parts {
                if strings.HasPrefix(part, ":") {
                    p.addParameter(strings.TrimPrefix(part, ":"), "path")
                }
            }
        }
    }
    
    result := make([]*Parameter, 0, len(p.parameters))
    for _, param := range p.parameters {
        result = append(result, param)
    }
    return result
}

// addParameter adds or updates a parameter
func (p *ParamMiner) addParameter(name, paramType string) {
    if existing, ok := p.parameters[name]; ok {
        if existing.Type != paramType {
            existing.Type = "multiple"
        }
    } else {
        p.parameters[name] = &Parameter{
            Name:     name,
            Type:     paramType,
            Examples: []string{},
            Required: false,
        }
    }
}

// getParamType determines parameter type from pattern
func (p *ParamMiner) getParamType(pattern string) string {
    if strings.Contains(pattern, "params") {
        return "path"
    }
    if strings.Contains(pattern, "query") || strings.Contains(pattern, "\\?") {
        return "query"
    }
    if strings.Contains(pattern, "body") || strings.Contains(pattern, "data") {
        return "body"
    }
    return "unknown"
}

// GetParameters returns all mined parameters
func (p *ParamMiner) GetParameters() []*Parameter {
    result := make([]*Parameter, 0, len(p.parameters))
    for _, param := range p.parameters {
        result = append(result, param)
    }
    return result
}

// GetQueryParams returns only query parameters
func (p *ParamMiner) GetQueryParams() []*Parameter {
    var queryParams []*Parameter
    for _, param := range p.parameters {
        if param.Type == "query" {
            queryParams = append(queryParams, param)
        }
    }
    return queryParams
}

// GetPathParams returns only path parameters
func (p *ParamMiner) GetPathParams() []*Parameter {
    var pathParams []*Parameter
    for _, param := range p.parameters {
        if param.Type == "path" {
            pathParams = append(pathParams, param)
        }
    }
    return pathParams
}
