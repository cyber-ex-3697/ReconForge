package portscan

import (
    "encoding/json"
    "os/exec"
    "strings"
)

type NaabuResult struct {
    Host  string `json:"host"`
    Port  int    `json:"port"`
    Name  string `json:"name"`
}

type NaabuWrapper struct {
    topPorts int
    threads  int
}

func NewNaabuWrapper(topPorts, threads int) *NaabuWrapper {
    return &NaabuWrapper{
        topPorts: topPorts,
        threads:  threads,
    }
}

func (n *NaabuWrapper) Scan(hosts []string) ([]NaabuResult, error) {
    if len(hosts) == 0 {
        return nil, nil
    }
    
    input := strings.Join(hosts, "\n")
    cmd := exec.Command("naabu", "-list", input, "-top-ports", string(rune(n.topPorts)), "-silent", "-json", "-c", string(rune(n.threads)))
    cmd.Stdin = strings.NewReader(input)
    output, err := cmd.Output()
    if err != nil {
        return nil, err
    }
    
    var results []NaabuResult
    for _, line := range strings.Split(string(output), "\n") {
        if line == "" {
            continue
        }
        var result NaabuResult
        if err := json.Unmarshal([]byte(line), &result); err == nil {
            results = append(results, result)
        }
    }
    return results, nil
}

func (n *NaabuWrapper) ScanSingle(host string) ([]NaabuResult, error) {
    return n.Scan([]string{host})
}

func (n *NaabuWrapper) GetOpenPorts(hosts []string) (map[string][]int, error) {
    results, err := n.Scan(hosts)
    if err != nil {
        return nil, err
    }
    
    ports := make(map[string][]int)
    for _, r := range results {
        ports[r.Host] = append(ports[r.Host], r.Port)
    }
    return ports, nil
}
