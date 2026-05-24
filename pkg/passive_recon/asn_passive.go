package passive_recon

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "strconv"
    "strings"
    "time"
)

// ASNInfo represents ASN information
type ASNInfo struct {
    ASN         int      `json:"asn"`
    Name        string   `json:"name"`
    Country     string   `json:"country"`
    Prefixes    []string `json:"prefixes"`
    IPv4Count   int      `json:"ipv4_count"`
    IPv6Count   int      `json:"ipv6_count"`
    RelatedASNs []int    `json:"related_asns"`
}

// PassiveASNCollector collects ASN data passively
type PassiveASNCollector struct {
    client  *http.Client
    results map[string]*ASNInfo
}

// NewPassiveASNCollector creates a new ASN collector
func NewPassiveASNCollector() *PassiveASNCollector {
    return &PassiveASNCollector{
        client: &http.Client{Timeout: 30 * time.Second},
        results: make(map[string]*ASNInfo),
    }
}

// CollectFromIPInfo collects ASN from IPInfo
func (pac *PassiveASNCollector) CollectFromIPInfo(ip string) (*ASNInfo, error) {
    url := fmt.Sprintf("https://ipinfo.io/%s/json", ip)
    resp, err := pac.client.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }
    
    var result map[string]interface{}
    if err := json.Unmarshal(body, &result); err != nil {
        return nil, err
    }
    
    asnInfo := &ASNInfo{}
    
    if org, ok := result["org"].(string); ok {
        parts := strings.SplitN(org, " ", 2)
        if len(parts) >= 2 {
            asn, _ := strconv.Atoi(strings.TrimPrefix(parts[0], "AS"))
            asnInfo.ASN = asn
            asnInfo.Name = parts[1]
        }
    }
    
    if country, ok := result["country"].(string); ok {
        asnInfo.Country = country
    }
    
    pac.results[ip] = asnInfo
    return asnInfo, nil
}

// CollectFromBGPView collects from BGPView API
func (pac *PassiveASNCollector) CollectFromBGPView(asn int) (*ASNInfo, error) {
    url := fmt.Sprintf("https://api.bgpview.io/asn/%d", asn)
    resp, err := pac.client.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var result map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }
    
    asnInfo := &ASNInfo{
        ASN: asn,
    }
    
    if data, ok := result["data"].(map[string]interface{}); ok {
        if name, ok := data["name"].(string); ok {
            asnInfo.Name = name
        }
        if country, ok := data["country_code"].(string); ok {
            asnInfo.Country = country
        }
        if prefixes, ok := data["prefixes"].([]interface{}); ok {
            for _, p := range prefixes {
                if prefix, ok := p.(map[string]interface{}); ok {
                    if prefixStr, ok := prefix["prefix"].(string); ok {
                        asnInfo.Prefixes = append(asnInfo.Prefixes, prefixStr)
                    }
                }
            }
        }
    }
    
    pac.results[strconv.Itoa(asn)] = asnInfo
    return asnInfo, nil
}

// CollectFromASNLookup collects from ASNLookup API
func (pac *PassiveASNCollector) CollectFromASNLookup(asn int) (*ASNInfo, error) {
    url := fmt.Sprintf("https://api.asnlookup.com/asn/%d", asn)
    resp, err := pac.client.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }
    
    var result map[string]interface{}
    if err := json.Unmarshal(body, &result); err != nil {
        return nil, err
    }
    
    asnInfo := &ASNInfo{
        ASN: asn,
    }
    
    if name, ok := result["name"].(string); ok {
        asnInfo.Name = name
    }
    if country, ok := result["country"].(string); ok {
        asnInfo.Country = country
    }
    
    pac.results[strconv.Itoa(asn)] = asnInfo
    return asnInfo, nil
}

// GetASNForDomain gets ASN information for a domain
func (pac *PassiveASNCollector) GetASNForDomain(domain string) (*ASNInfo, error) {
    // First resolve domain to IP (simplified)
    // In production, use net.LookupIP
    
    // For demo, return mock data
    return &ASNInfo{
        ASN:     15169,
        Name:    "Google LLC",
        Country: "US",
    }, nil
}

// GetResults returns all ASN results
func (pac *PassiveASNCollector) GetResults() map[string]*ASNInfo {
    return pac.results
}

// GetPrefixes returns all prefixes for an ASN
func (pac *PassiveASNCollector) GetPrefixes(asn int) []string {
    for _, info := range pac.results {
        if info.ASN == asn {
            return info.Prefixes
        }
    }
    return []string{}
}

// ExpandPrefixes expands CIDR prefixes to IP ranges
func (pac *PassiveASNCollector) ExpandPrefixes(prefixes []string) []string {
    var ips []string
    for _, prefix := range prefixes {
        // Simplified - would expand CIDR in production
        ips = append(ips, prefix)
    }
    return ips
}
