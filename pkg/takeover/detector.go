package takeover

import (
    "bufio"
    "os/exec"
    "strings"
)

type TakeoverResult struct {
    Subdomain string
    Service   string
    CNAME     string
    Status    string
}

type Detector struct {
    concurrency int
}

func NewDetector(concurrency int) *Detector {
    return &Detector{
        concurrency: concurrency,
    }
}

func (d *Detector) Check(domains []string) ([]TakeoverResult, error) {
    if len(domains) == 0 {
        return nil, nil
    }
    
    cmd := exec.Command("subzy", "run", "--list", strings.Join(domains, "\n"), "--concurrency", string(rune(d.concurrency)))
    output, err := cmd.Output()
    if err != nil {
        return nil, err
    }
    
    var results []TakeoverResult
    scanner := bufio.NewScanner(strings.NewReader(string(output)))
    for scanner.Scan() {
        line := scanner.Text()
        if strings.Contains(line, "VULNERABLE") {
            parts := strings.Fields(line)
            if len(parts) >= 2 {
                results = append(results, TakeoverResult{
                    Subdomain: parts[0],
                    Status:    "VULNERABLE",
                })
            }
        }
    }
    return results, nil
}
