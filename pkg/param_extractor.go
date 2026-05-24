package url

import (
    "strings"
)

type ParamExtractor struct {
    seen map[string]bool
}

func NewParamExtractor() *ParamExtractor {
    return &ParamExtractor{
        seen: make(map[string]bool),
    }
}

func (p *ParamExtractor) ExtractParams(url string) []string {
    var params []string
    
    if !strings.Contains(url, "?") {
        return params
    }
    
    queryPart := strings.Split(url, "?")[1]
    pairs := strings.Split(queryPart, "&")
    
    for _, pair := range pairs {
        if strings.Contains(pair, "=") {
            param := strings.Split(pair, "=")[0]
            if !p.seen[param] {
                p.seen[param] = true
                params = append(params, param)
            }
        }
    }
    return params
}

func (p *ParamExtractor) ExtractParamsBatch(urls []string) []string {
    var allParams []string
    for _, url := range urls {
        params := p.ExtractParams(url)
        allParams = append(allParams, params...)
    }
    return allParams
}

func (p *ParamExtractor) GetUniqueParams(urls []string) []string {
    params := p.ExtractParamsBatch(urls)
    seen := make(map[string]bool)
    var unique []string
    
    for _, param := range params {
        if !seen[param] {
            seen[param] = true
            unique = append(unique, param)
        }
    }
    return unique
}

func (p *ParamExtractor) Reset() {
    p.seen = make(map[string]bool)
}

func (p *ParamExtractor) GetParamCount() int {
    return len(p.seen)
}
