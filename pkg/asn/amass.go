package asn

import (
    "bufio"
    "os/exec"
    "strings"
)

type AmassIntel struct {
    timeout int
}

func NewAmassIntel(timeout int) *AmassIntel {
    return &AmassIntel{
        timeout: timeout,
    }
}

func (a *AmassIntel) GetASN(domain string) ([]ASNInfo, error) {
    cmd := exec.Command("amass", "intel", "-whois", "-d", domain)
    output, err := cmd.Output()
    if err != nil {
        return nil, err
    }
    
    var results []ASNInfo
    scanner := bufio.NewScanner(strings.NewReader(string(output)))
    for scanner.Scan() {
        line := scanner.Text()
        if strings.Contains(line, "AS") {
            parts := strings.Fields(line)
            if len(parts) >= 2 {
                var asn int
                // Parse ASN number from ASxxxx
                asnStr := strings.TrimPrefix(parts[1], "AS")
                // Manual parsing
                results = append(results, ASNInfo{
                    IP:  parts[0],
                    ASN: asn,
                    Org: strings.Join(parts[2:], " "),
                })
            }
        }
    }
    return results, nil
}

func (a *AmassIntel) GetIPRanges(domain string) ([]string, error) {
    cmd := exec.Command("amass", "intel", "-whois", "-d", domain, "-cidr")
    output, err := cmd.Output()
    if err != nil {
        return nil, err
    }
    
    var ranges []string
    scanner := bufio.NewScanner(strings.NewReader(string(output)))
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if line != "" {
            ranges = append(ranges, line)
        }
    }
    return ranges, nil
}

func (a *AmassIntel) GetOrganizations(domain string) ([]string, error) {
    cmd := exec.Command("amass", "intel", "-whois", "-d", domain, "-org")
    output, err := cmd.Output()
    if err != nil {
        return nil, err
    }
    
    var orgs []string
    scanner := bufio.NewScanner(strings.NewReader(string(output)))
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if line != "" {
            orgs = append(orgs, line)
        }
    }
    return orgs, nil
}
