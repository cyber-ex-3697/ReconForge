package portscan

import (
    "bufio"
    "os/exec"
    "strconv"
    "strings"
)

type MasscanResult struct {
    IP       string
    Port     int
    Protocol string
    Status   string
}

type MasscanWrapper struct {
    rate    int
    ports   string
}

func NewMasscanWrapper(rate int) *MasscanWrapper {
    return &MasscanWrapper{
        rate:  rate,
        ports: "0-65535",
    }
}

func (m *MasscanWrapper) SetPorts(ports string) {
    m.ports = ports
}

func (m *MasscanWrapper) Scan(target string) ([]MasscanResult, error) {
    cmd := exec.Command("masscan", target, "-p", m.ports, "--rate", strconv.Itoa(m.rate), "-oG", "-")
    output, err := cmd.Output()
    if err != nil {
        return nil, err
    }
    
    var results []MasscanResult
    scanner := bufio.NewScanner(strings.NewReader(string(output)))
    
    for scanner.Scan() {
        line := scanner.Text()
        if strings.Contains(line, "open") {
            parts := strings.Fields(line)
            for i, part := range parts {
                if part == "open" && i+1 < len(parts) {
                    portProto := strings.Split(parts[i+1], "/")
                    if len(portProto) >= 2 {
                        port, _ := strconv.Atoi(portProto[0])
                        results = append(results, MasscanResult{
                            IP:       parts[1],
                            Port:     port,
                            Protocol: portProto[1],
                            Status:   "open",
                        })
                    }
                }
            }
        }
    }
    return results, nil
}

func (m *MasscanWrapper) ScanCIDR(cidr string) ([]MasscanResult, error) {
    return m.Scan(cidr)
}

func (m *MasscanWrapper) GetOpenPorts(target string) ([]int, error) {
    results, err := m.Scan(target)
    if err != nil {
        return nil, err
    }
    
    var ports []int
    for _, r := range results {
        ports = append(ports, r.Port)
    }
    return ports, nil
}
