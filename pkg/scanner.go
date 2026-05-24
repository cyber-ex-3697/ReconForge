package vulnerability

import (
    "bufio"
    "encoding/json"
    "fmt"
    "os"
    "os/exec"
    "strings"
)

type Scanner struct {
    rateLimit int
    severity  string
}

type Vulnerability struct {
    URL        string  `json:"url"`
    TemplateID string  `json:"template_id"`
    Severity   string  `json:"severity"`
    Name       string  `json:"name"`
    MatchedAt  string  `json:"matched_at"`
    CVSS       float64 `json:"cvss"`
}

func NewScanner(rateLimit int) *Scanner {
    return &Scanner{
        rateLimit: rateLimit,
        severity:  "critical,high,medium",
    }
}

func (s *Scanner) SetSeverity(severity string) {
    s.severity = severity
}

func (s *Scanner) Scan(hosts []string) ([]Vulnerability, error) {
    if len(hosts) == 0 {
        return nil, nil
    }
    
    tempFile := "/tmp/nuclei_hosts.txt"
    content := strings.Join(hosts, "\n")
    os.WriteFile(tempFile, []byte(content), 0644)
    
    cmd := exec.Command("nuclei", "-l", tempFile, "-silent", "-severity", s.severity, "-rate-limit", fmt.Sprintf("%d", s.rateLimit), "-json")
    output, err := cmd.Output()
    if err != nil {
        return nil, err
    }
    
    var vulnerabilities []Vulnerability
    for _, line := range strings.Split(string(output), "\n") {
        if line == "" {
            continue
        }
        var vuln Vulnerability
        if err := json.Unmarshal([]byte(line), &vuln); err == nil {
            vulnerabilities = append(vulnerabilities, vuln)
        }
    }
    return vulnerabilities, nil
}

func (s *Scanner) GetCriticalVulnerabilities(hosts []string) ([]Vulnerability, error) {
    s.severity = "critical"
    return s.Scan(hosts)
}

func (s *Scanner) GetHighVulnerabilities(hosts []string) ([]Vulnerability, error) {
    s.severity = "high"
    return s.Scan(hosts)
}
