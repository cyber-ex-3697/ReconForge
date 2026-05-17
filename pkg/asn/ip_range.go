package asn

import (
    "net"
    "strconv"
)

type IPRange struct {
    Start net.IP
    End   net.IP
    CIDR  string
}

func (i *IPRange) Contains(ip net.IP) bool {
    return bytesCompare(i.Start, ip) <= 0 && bytesCompare(ip, i.End) <= 0
}

func bytesCompare(a, b net.IP) int {
    for i := 0; i < len(a) && i < len(b); i++ {
        if a[i] < b[i] {
            return -1
        }
        if a[i] > b[i] {
            return 1
        }
    }
    return 0
}

func ParseCIDR(cidr string) (*IPRange, error) {
    _, ipnet, err := net.ParseCIDR(cidr)
    if err != nil {
        return nil, err
    }
    
    start := ipnet.IP
    end := make(net.IP, len(start))
    copy(end, start)
    
    // Calculate end IP
    for i := range end {
        end[i] |= ^ipnet.Mask[i]
    }
    
    return &IPRange{
        Start: start,
        End:   end,
        CIDR:  cidr,
    }, nil
}

func GetIPRangeFromASN(asn int) ([]string, error) {
    // This would require a prefix database
    // For now, return empty
    return []string{}, nil
}

func ExpandCIDR(cidr string) ([]string, error) {
    _, ipnet, err := net.ParseCIDR(cidr)
    if err != nil {
        return nil, err
    }
    
    var ips []string
    ip := ipnet.IP.Mask(ipnet.Mask)
    for {
        ips = append(ips, ip.String())
        for i := len(ip) - 1; i >= 0; i-- {
            ip[i]++
            if ip[i] != 0 {
                break
            }
        }
        if !ipnet.Contains(ip) {
            break
        }
    }
    return ips, nil
}
