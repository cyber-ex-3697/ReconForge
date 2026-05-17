package subdomain

import (
    "strings"
)

type RecursiveEnumerator struct {
    target    string
    maxDepth  int
    discovered map[string]bool
}

func NewRecursiveEnumerator(target string, maxDepth int) *RecursiveEnumerator {
    return &RecursiveEnumerator{
        target:     target,
        maxDepth:   maxDepth,
        discovered: make(map[string]bool),
    }
}

func (r *RecursiveEnumerator) Enumerate(depth int) ([]string, error) {
    if depth > r.maxDepth {
        result := make([]string, 0, len(r.discovered))
        for s := range r.discovered {
            result = append(result, s)
        }
        return result, nil
    }
    
    // Get current level subdomains
    enumerator := NewEnumerator(r.target, 50)
    result, err := enumerator.Run()
    if err != nil {
        return nil, err
    }
    
    for _, sub := range result.Subdomains {
        r.discovered[sub] = true
    }
    
    // Recursively enumerate each subdomain
    for _, sub := range result.Subdomains {
        parts := strings.Split(sub, ".")
        if len(parts) > 2 {
            newTarget := strings.Join(parts[1:], ".")
            r.target = newTarget
            return r.Enumerate(depth + 1)
        }
    }
    
    finalResult := make([]string, 0, len(r.discovered))
    for s := range r.discovered {
        finalResult = append(finalResult, s)
    }
    return finalResult, nil
}
