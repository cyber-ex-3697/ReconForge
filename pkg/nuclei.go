package vulnerability

import (
    "encoding/json"
    "fmt"
    "os/exec"
    "strings"
)

type NucleiResult struct {
    TemplateID string  `json:"template-id"`
    Name       string  `json:"name"`
    Severity   string  `json:"severity"`
    MatchedAt  string  `json:"matched-at"`
    CVSS       float64 `json:"cvss"`
    Description string `json:"description"`
}

type NucleiWrapper struct {
    rateLimit  int
    severity   string
    concurrency int
}

func NewNucleiWrapper(rateLimit int) *NucleiWrapper {
    return &NucleiWrapper{
        rateLimit:  rateLimit,
        severity:   "critical,high,medium",
        concurrency: 20,
    }
}

func (n *NucleiWrapper) SetSeverity(severity string) {
    n.severity = severity
}

func (n *NucleiWrapper) SetConcurrency(conc int) {
    n.concurrency = conc
}

func (n *NucleiWrapper) Scan(targets []string) ([]NucleiResult, error) {
    if len(targets) == 0 {
        return nil, nil
    }
    
    input := strings.Join(targets, "\n")
    cmd := exec.Command("nuclei", "-silent", "-severity", n.severity, 
        "-rate-limit", fmt.Sprintf("%d", n.rateLimit), 
        "-concurrency", fmt.Sprintf("%d", n.concurrency),
        "-json")
    cmd.Stdin = strings.NewReader(input)
    
    output, err := cmd.Output()
    if err != nil {
        return nil, err
    }
    
    var results []NucleiResult
    for _, line := range strings.Split(string(output), "\n") {
        if line == "" {
            continue
        }
        var result NucleiResult
        if err := json.Unmarshal([]byte(line), &result); err == nil {
            results = append(results, result)
        }
    }
    return results, nil
}

func (n *NucleiWrapper) ScanWithTemplates(targets []string, templates []string) ([]NucleiResult, error) {
    if len(targets) == 0 {
        return nil, nil
    }
    
    templateStr := strings.Join(templates, ",")
    input := strings.Join(targets, "\n")
    cmd := exec.Command("nuclei", "-silent", "-t", templateStr,
        "-rate-limit", fmt.Sprintf("%d", n.rateLimit),
        "-json")
    cmd.Stdin = strings.NewReader(input)
    
    output, err := cmd.Output()
    if err != nil {
        return nil, err
    }
    
    var results []NucleiResult
    for _, line := range strings.Split(string(output), "\n") {
        if line == "" {
            continue
        }
        var result NucleiResult
        if err := json.Unmarshal([]byte(line), &result); err == nil {
            results = append(results, result)
        }
    }
    return results, nil
}

func (n *NucleiWrapper) GetCriticalFindings(targets []string) ([]NucleiResult, error) {
    n.severity = "critical"
    return n.Scan(targets)
}

func (n *NucleiWrapper) GetHighFindings(targets []string) ([]NucleiResult, error) {
    n.severity = "high"
    return n.Scan(targets)
}

func (n *NucleiWrapper) GetMediumFindings(targets []string) ([]NucleiResult, error) {
    n.severity = "medium"
    return n.Scan(targets)
}
