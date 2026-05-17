package portscan

import (
    "bufio"
    "fmt"
    "os"
    "os/exec"
    "strings"
)

type Scanner struct {
    threads  int
    timeout  int
    ports    string
}

type PortResult struct {
    Host      string
    Port      int
    Service   string
    Protocol  string
}

func NewScanner(threads int) *Scanner {
    return &Scanner{
        threads: threads,
        timeout: 10,
        ports:   "80,443,22,21,25,53,110,143,993,995,3306,5432,8080,8443,8000,8888,27017,6379,9200,5000",
    }
}

func (s *Scanner) SetPorts(ports string) {
    s.ports = ports
}

func (s *Scanner) ScanWithNaabu(hosts []string) ([]PortResult, error) {
    if len(hosts) == 0 {
        return nil, nil
    }
    
    tempFile := "/tmp/naabu_hosts.txt"
    content := strings.Join(hosts, "\n")
    os.WriteFile(tempFile, []byte(content), 0644)
    
    cmd := exec.Command("naabu", "-list", tempFile, "-top-ports", "1000", "-silent", "-json")
    output, err := cmd.Output()
    if err != nil {
        return nil, err
    }
    
    var results []PortResult
    for _, line := range strings.Split(string(output), "\n") {
        if line == "" {
            continue
        }
        // Parse JSON output
        var result PortResult
        // Simple parsing for now
        if strings.Contains(line, "port") {
            results = append(results, result)
        }
    }
    return results, nil
}

func (s *Scanner) ScanWithNmap(host string) ([]PortResult, error) {
    cmd := exec.Command("nmap", "-p", s.ports, "--open", "-T4", host)
    output, err := cmd.Output()
    if err != nil {
        return nil, err
    }
    
    var results []PortResult
    scanner := bufio.NewScanner(strings.NewReader(string(output)))
    for scanner.Scan() {
        line := scanner.Text()
        if strings.Contains(line, "/tcp") || strings.Contains(line, "/udp") {
            parts := strings.Fields(line)
            if len(parts) >= 3 {
                portParts := strings.Split(parts[0], "/")
                var port int
                fmt.Sscanf(portParts[0], "%d", &port)
                results = append(results, PortResult{
                    Host:     host,
                    Port:     port,
                    Protocol: portParts[1],
                    Service:  parts[2],
                })
            }
        }
    }
    return results, nil
}
