package subdomain

import (
    "os/exec"
    "strings"
)

type PassiveSource struct {
    name   string
    target string
}

func NewPassiveSource(target string) *PassiveSource {
    return &PassiveSource{
        target: target,
    }
}

func (p *PassiveSource) RunSubfinder() ([]string, error) {
    cmd := exec.Command("subfinder", "-d", p.target, "-silent")
    output, err := cmd.Output()
    if err != nil {
        return nil, err
    }
    
    var results []string
    for _, s := range strings.Split(string(output), "\n") {
        if s != "" {
            results = append(results, s)
        }
    }
    return results, nil
}

func (p *PassiveSource) RunAssetfinder() ([]string, error) {
    cmd := exec.Command("assetfinder", "--subs-only", p.target)
    output, err := cmd.Output()
    if err != nil {
        return nil, err
    }
    
    var results []string
    for _, s := range strings.Split(string(output), "\n") {
        if s != "" {
            results = append(results, s)
        }
    }
    return results, nil
}

func (p *PassiveSource) RunFindomain() ([]string, error) {
    cmd := exec.Command("findomain", "-t", p.target, "-q")
    output, err := cmd.Output()
    if err != nil {
        return nil, err
    }
    
    var results []string
    for _, s := range strings.Split(string(output), "\n") {
        if s != "" {
            results = append(results, s)
        }
    }
    return results, nil
}

func (p *PassiveSource) RunAll() ([]string, error) {
    all := make(map[string]bool)
    
    if subs, err := p.RunSubfinder(); err == nil {
        for _, s := range subs {
            all[s] = true
        }
    }
    
    if subs, err := p.RunAssetfinder(); err == nil {
        for _, s := range subs {
            all[s] = true
        }
    }
    
    if subs, err := p.RunFindomain(); err == nil {
        for _, s := range subs {
            all[s] = true
        }
    }
    
    result := make([]string, 0, len(all))
    for s := range all {
        result = append(result, s)
    }
    return result, nil
}
