package asn

import (
    "bufio"
    "os/exec"
    "strconv"
    "strings"
)

type ASNInfo struct {
    IP      string
    ASN     int
    Org     string
    Country string
    Prefix  string
}

type Enumerator struct {
    timeout int
}

func NewEnumerator(timeout int) *Enumerator {
    return &Enumerator{
        timeout: timeout,
    }
}

func (e *Enumerator) GetASN(ip string) (*ASNInfo, error) {
    cmd := exec.Command("asnmap", "-ip", ip, "-silent")
    output, err := cmd.Output()
    if err != nil {
        return nil, err
    }
    
    parts := strings.Fields(string(output))
    if len(parts) < 5 {
        return nil, nil
    }
    
    asn, _ := strconv.Atoi(strings.TrimPrefix(parts[1], "AS"))
    
    return &ASNInfo{
        IP:      parts[0],
        ASN:     asn,
        Org:     parts[2],
        Country: parts[3],
        Prefix:  parts[4],
    }, nil
}

func (e *Enumerator) GetASNBatch(ips []string) ([]ASNInfo, error) {
    var results []ASNInfo
    for _, ip := range ips {
        info, err := e.GetASN(ip)
        if err == nil && info != nil {
            results = append(results, *info)
        }
    }
    return results, nil
}

func (e *Enumerator) GetASNFromCIDR(cidr string) ([]ASNInfo, error) {
    cmd := exec.Command("asnmap", "-cidr", cidr, "-silent")
    output, err := cmd.Output()
    if err != nil {
        return nil, err
    }
    
    var results []ASNInfo
    scanner := bufio.NewScanner(strings.NewReader(string(output)))
    for scanner.Scan() {
        line := scanner.Text()
        parts := strings.Fields(line)
        if len(parts) >= 5 {
            asn, _ := strconv.Atoi(strings.TrimPrefix(parts[1], "AS"))
            results = append(results, ASNInfo{
                IP:      parts[0],
                ASN:     asn,
                Org:     parts[2],
                Country: parts[3],
                Prefix:  parts[4],
            })
        }
    }
    return results, nil
}
