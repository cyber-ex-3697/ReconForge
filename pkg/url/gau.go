package url

import (
    "os/exec"
    "strings"
)

type GAUWrapper struct {
    subs bool
}

func NewGAUWrapper() *GAUWrapper {
    return &GAUWrapper{
        subs: true,
    }
}

func (g *GAUWrapper) GetURLs(domain string) ([]string, error) {
    var cmd *exec.Cmd
    if g.subs {
        cmd = exec.Command("gau", "--subs", domain)
    } else {
        cmd = exec.Command("gau", domain)
    }
    
    output, err := cmd.Output()
    if err != nil {
        return nil, err
    }
    
    var urls []string
    for _, u := range strings.Split(string(output), "\n") {
        if u != "" {
            urls = append(urls, u)
        }
    }
    return urls, nil
}

func (g *GAUWrapper) SetSubs(enabled bool) {
    g.subs = enabled
}

func (g *GAUWrapper) GetUniqueURLs(domain string) ([]string, error) {
    urls, err := g.GetURLs(domain)
    if err != nil {
        return nil, err
    }
    
    seen := make(map[string]bool)
    var unique []string
    for _, u := range urls {
        if !seen[u] {
            seen[u] = true
            unique = append(unique, u)
        }
    }
    return unique, nil
}
