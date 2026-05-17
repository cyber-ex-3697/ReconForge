package utils

import (
    "net"
    "strings"
    "time"
)

type Resolver struct {
    timeout time.Duration
}

func NewResolver() *Resolver {
    return &Resolver{
        timeout: 5 * time.Second,
    }
}

func (r *Resolver) Resolve(domain string) ([]string, error) {
    ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
    defer cancel()
    
    var resolver net.Resolver
    ips, err := resolver.LookupHost(ctx, domain)
    if err != nil {
        return nil, err
    }
    
    return ips, nil
}

func (r *Resolver) ResolveBatch(domains []string) map[string][]string {
    results := make(map[string][]string)
    
    for _, domain := range domains {
        ips, err := r.Resolve(domain)
        if err == nil && len(ips) > 0 {
            results[domain] = ips
        }
    }
    
    return results
}

func (r *Resolver) HasCNAME(domain string) (bool, string, error) {
    cname, err := net.LookupCNAME(domain)
    if err != nil {
        return false, "", err
    }
    return true, strings.TrimSuffix(cname, "."), nil
}
