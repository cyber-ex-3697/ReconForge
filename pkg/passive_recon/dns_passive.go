package passive_recon

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "strings"
    "time"
)

// PassiveDNSCollector collects DNS data from passive sources
type PassiveDNSCollector struct {
    client    *http.Client
    sources   []string
    results   map[string][]string
}

// DNSRecord represents a DNS record from passive sources
type DNSRecord struct {
    Domain    string    `json:"domain"`
    Type      string    `json:"type"`
    Value     string    `json:"value"`
    FirstSeen time.Time `json:"first_seen"`
    LastSeen  time.Time `json:"last_seen"`
    Source    string    `json:"source"`
}

// NewPassiveDNSCollector creates a new passive DNS collector
func NewPassiveDNSCollector() *PassiveDNSCollector {
    return &PassiveDNSCollector{
        client: &http.Client{Timeout: 30 * time.Second},
        sources: []string{
            "https://api.securitytrails.com/v1/domain/{domain}/subdomains",
            "https://dns.bufferover.run/dns?q=.{domain}",
            "https://api.omnisint.io/subdomains/{domain}",
        },
        results: make(map[string][]string),
    }
}

// CollectFromSecurityTrails collects from SecurityTrails API
func (pdc *PassiveDNSCollector) CollectFromSecurityTrails(domain, apiKey string) ([]string, error) {
    url := fmt.Sprintf("https://api.securitytrails.com/v1/domain/%s/subdomains", domain)
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, err
    }
    req.Header.Set("APIKEY", apiKey)
    
    resp, err := pdc.client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }
    
    var result struct {
        Subdomains []string `json:"subdomains"`
    }
    if err := json.Unmarshal(body, &result); err != nil {
        return nil, err
    }
    
    var domains []string
    for _, sub := range result.Subdomains {
        domains = append(domains, fmt.Sprintf("%s.%s", sub, domain))
    }
    
    pdc.results[domain] = append(pdc.results[domain], domains...)
    return domains, nil
}

// CollectFromBufferOver collects from BufferOver.run
func (pdc *PassiveDNSCollector) CollectFromBufferOver(domain string) ([]string, error) {
    url := fmt.Sprintf("https://dns.bufferover.run/dns?q=.%s", domain)
    resp, err := pdc.client.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }
    
    var result struct {
        FDNSA []string `json:"FDNS_A"`
    }
    if err := json.Unmarshal(body, &result); err != nil {
        return nil, err
    }
    
    var domains []string
    for _, record := range result.FDNSA {
        parts := strings.Split(record, ",")
        if len(parts) >= 2 {
            domains = append(domains, parts[1])
        }
    }
    
    pdc.results[domain] = append(pdc.results[domain], domains...)
    return domains, nil
}

// CollectFromOmnisint collects from Omnisint API
func (pdc *PassiveDNSCollector) CollectFromOmnisint(domain string) ([]string, error) {
    url := fmt.Sprintf("https://api.omnisint.io/subdomains/%s", domain)
    resp, err := pdc.client.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var subdomains []string
    if err := json.NewDecoder(resp.Body).Decode(&subdomains); err != nil {
        return nil, err
    }
    
    var domains []string
    for _, sub := range subdomains {
        domains = append(domains, fmt.Sprintf("%s.%s", sub, domain))
    }
    
    pdc.results[domain] = append(pdc.results[domain], domains...)
    return domains, nil
}

// CollectAll collects from all available passive sources
func (pdc *PassiveDNSCollector) CollectAll(domain string, apiKey string) ([]string, error) {
    allDomains := make(map[string]bool)
    
    // BufferOver
    if domains, err := pdc.CollectFromBufferOver(domain); err == nil {
        for _, d := range domains {
            allDomains[d] = true
        }
    }
    
    // Omnisint
    if domains, err := pdc.CollectFromOmnisint(domain); err == nil {
        for _, d := range domains {
            allDomains[d] = true
        }
    }
    
    // SecurityTrails (requires API key)
    if apiKey != "" {
        if domains, err := pdc.CollectFromSecurityTrails(domain, apiKey); err == nil {
            for _, d := range domains {
                allDomains[d] = true
            }
        }
    }
    
    result := make([]string, 0, len(allDomains))
    for d := range allDomains {
        result = append(result, d)
    }
    
    return result, nil
}

// GetResults returns all collected domains
func (pdc *PassiveDNSCollector) GetResults() map[string][]string {
    return pdc.results
}

// Clear clears collected results
func (pdc *PassiveDNSCollector) Clear() {
    pdc.results = make(map[string][]string)
}
