package takeover

import (
    "bufio"
    "encoding/json"
    "os/exec"
    "strings"
)

type SubzyResult struct {
    Hostname  string `json:"hostname"`
    Service   string `json:"service"`
    Vulnerable bool  `json:"vulnerable"`
    CNAME     string `json:"cname"`
}

type SubzyWrapper struct {
    concurrency int
    timeout     int
}

func NewSubzyWrapper(concurrency, timeout int) *SubzyWrapper {
    return &SubzyWrapper{
        concurrency: concurrency,
        timeout:     timeout,
    }
}

func (s *SubzyWrapper) Check(domains []string) ([]SubzyResult, error) {
    if len(domains) == 0 {
        return nil, nil
    }
    
    input := strings.Join(domains, "\n")
    cmd := exec.Command("subzy", "run", "--list", input, 
        "--concurrency", string(rune(s.concurrency)),
        "--timeout", string(rune(s.timeout)),
        "--json")
    cmd.Stdin = strings.NewReader(input)
    
    output, err := cmd.Output()
    if err != nil {
        return nil, err
    }
    
    var results []SubzyResult
    for _, line := range strings.Split(string(output), "\n") {
        if line == "" {
            continue
        }
        var result SubzyResult
        if err := json.Unmarshal([]byte(line), &result); err == nil {
            results = append(results, result)
        }
    }
    return results, nil
}

func (s *SubzyWrapper) GetVulnerable(domains []string) ([]SubzyResult, error) {
    results, err := s.Check(domains)
    if err != nil {
        return nil, err
    }
    
    var vulnerable []SubzyResult
    for _, r := range results {
        if r.Vulnerable {
            vulnerable = append(vulnerable, r)
        }
    }
    return vulnerable, nil
}

func (s *SubzyWrapper) CheckSingle(domain string) (*SubzyResult, error) {
    results, err := s.Check([]string{domain})
    if err != nil {
        return nil, err
    }
    if len(results) > 0 {
        return &results[0], nil
    }
    return nil, nil
}
