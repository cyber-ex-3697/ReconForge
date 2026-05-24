package subdomain

import (
    "os/exec"
    "strings"
)

type Resolver struct {
    resolvers []string
}

func NewResolver() *Resolver {
    return &Resolver{
        resolvers: []string{"1.1.1.1", "8.8.8.8"},
    }
}

func (r *Resolver) Resolve(domains []string) (map[string][]string, error) {
    input := strings.Join(domains, "\n")
    cmd := exec.Command("dnsx", "-silent", "-a", "-resp")
    cmd.Stdin = strings.NewReader(input)
    output, err := cmd.Output()
    if err != nil {
        return nil, err
    }
    
    results := make(map[string][]string)
    for _, line := range strings.Split(string(output), "\n") {
        if line == "" {
            continue
        }
        parts := strings.Fields(line)
        if len(parts) >= 2 {
            domain := parts[0]
            ips := parts[1:]
            results[domain] = ips
        }
    }
    return results, nil
}

func (r *Resolver) ResolveBatch(domains []string, threads int) (map[string][]string, error) {
    return r.Resolve(domains)
}

func (r *Resolver) GetResolvedDomains(domains []string) ([]string, error) {
    resolved, err := r.Resolve(domains)
    if err != nil {
        return nil, err
    }
    
    result := make([]string, 0, len(resolved))
    for domain := range resolved {
        result = append(result, domain)
    }
    return result, nil
}
