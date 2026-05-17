package livehost

import (
    "bufio"
    "fmt"
    "os"
    "os/exec"
    "strings"
)

type Prober struct {
    threads int
    timeout int
}

type LiveHostResult struct {
    URL        string
    StatusCode int
    Title      string
    TechStack  []string
}

func NewProber(threads int) *Prober {
    return &Prober{
        threads: threads,
        timeout: 10,
    }
}

func (p *Prober) Probe(domains []string) ([]LiveHostResult, error) {
    if len(domains) == 0 {
        return nil, nil
    }
    
    // Create temp file
    tempFile := "/tmp/probe_domains.txt"
    content := strings.Join(domains, "\n")
    os.WriteFile(tempFile, []byte(content), 0644)
    
    cmd := exec.Command("httpx", "-l", tempFile, "-silent", "-status-code", "-title", "-tech-detect", "-threads", fmt.Sprintf("%d", p.threads))
    output, err := cmd.Output()
    if err != nil {
        return nil, err
    }
    
    var results []LiveHostResult
    scanner := bufio.NewScanner(strings.NewReader(string(output)))
    for scanner.Scan() {
        line := scanner.Text()
        if strings.Contains(line, "[200]") || strings.Contains(line, "[301]") || strings.Contains(line, "[302]") {
            parts := strings.Fields(line)
            if len(parts) > 0 {
                result := LiveHostResult{
                    URL: parts[0],
                }
                results = append(results, result)
            }
        }
    }
    return results, nil
}

func (p *Prober) GetLiveURLs(domains []string) ([]string, error) {
    results, err := p.Probe(domains)
    if err != nil {
        return nil, err
    }
    
    var urls []string
    for _, r := range results {
        urls = append(urls, r.URL)
    }
    return urls, nil
}
