package portscan

import (
    "bufio"
    "os/exec"
    "strconv"
    "strings"
)

type NmapResult struct {
    Host      string
    Port      int
    Protocol  string
    State     string
    Service   string
    Version   string
}

type NmapWrapper struct {
    ports string
    fast  bool
}

func NewNmapWrapper() *NmapWrapper {
    return &NmapWrapper{
        ports: "80,443,22,21,25,53,110,143,993,995,3306,5432,8080,8443,8000,8888",
        fast:  true,
    }
}

func (n *NmapWrapper) SetPorts(ports string) {
    n.ports = ports
}

func (n *NmapWrapper) Scan(host string) ([]NmapResult, error) {
    var cmd *exec.Cmd
    if n.fast {
        cmd = exec.Command("nmap", "-p", n.ports, "--open", "-T4", host)
    } else {
        cmd = exec.Command("nmap", "-p", n.ports, "-sV", "-sC", "-T4", host)
    }
    
    output, err := cmd.Output()
    if err != nil {
        return nil, err
    }
    
    var results []NmapResult
    scanner := bufio.NewScanner(strings.NewReader(string(output)))
    
    for scanner.Scan() {
        line := scanner.Text()
        if strings.Contains(line, "/tcp") || strings.Contains(line, "/udp") {
            parts := strings.Fields(line)
            if len(parts) >= 3 {
                portProto := strings.Split(parts[0], "/")
                if len(portProto) >= 2 {
                    port, _ := strconv.Atoi(portProto[0])
                    results = append(results, NmapResult{
                        Host:     host,
                        Port:     port,
                        Protocol: portProto[1],
                        State:    parts[1],
                        Service:  parts[2],
                    })
                }
            }
        }
    }
    return results, nil
}

func (n *NmapWrapper) ScanBatch(hosts []string) (map[string][]NmapResult, error) {
    results := make(map[string][]NmapResult)
    for _, host := range hosts {
        hostResults, err := n.Scan(host)
        if err == nil {
            results[host] = hostResults
        }
    }
    return results, nil
}
