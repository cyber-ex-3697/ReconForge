package passive_recon

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "strings"
    "time"
)

// CertTransparencyCollector collects subdomains from certificate transparency logs
type CertTransparencyCollector struct {
    client  *http.Client
    results map[string][]string
}

// CertificateInfo represents certificate information
type CertificateInfo struct {
    Issuer      string    `json:"issuer"`
    Subject     string    `json:"subject"`
    NotBefore   time.Time `json:"not_before"`
    NotAfter    time.Time `json:"not_after"`
    Domains     []string  `json:"domains"`
}

// NewCertTransparencyCollector creates a new CT collector
func NewCertTransparencyCollector() *CertTransparencyCollector {
    return &CertTransparencyCollector{
        client: &http.Client{Timeout: 30 * time.Second},
        results: make(map[string][]string),
    }
}

// CollectFromCrtSh collects from crt.sh
func (ctc *CertTransparencyCollector) CollectFromCrtSh(domain string) ([]string, error) {
    url := fmt.Sprintf("https://crt.sh/?q=%%25.%s&output=json", domain)
    resp, err := ctc.client.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }
    
    var certificates []map[string]interface{}
    if err := json.Unmarshal(body, &certificates); err != nil {
        return nil, err
    }
    
    domains := make(map[string]bool)
    for _, cert := range certificates {
        if name, ok := cert["name_value"].(string); ok {
            for _, d := range strings.Split(name, "\n") {
                d = strings.TrimSpace(d)
                if strings.Contains(d, domain) && d != "" {
                    // Clean up domain (remove wildcard)
                    d = strings.TrimPrefix(d, "*.")
                    domains[d] = true
                }
            }
        }
    }
    
    result := make([]string, 0, len(domains))
    for d := range domains {
        result = append(result, d)
    }
    
    ctc.results[domain] = append(ctc.results[domain], result...)
    return result, nil
}

// CollectFromCertSpotter collects from CertSpotter
func (ctc *CertTransparencyCollector) CollectFromCertSpotter(domain string) ([]string, error) {
    url := fmt.Sprintf("https://api.certspotter.com/v1/issuances?domain=%s&include_subdomains=true&expand=dns_names", domain)
    resp, err := ctc.client.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var issuances []map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&issuances); err != nil {
        return nil, err
    }
    
    domains := make(map[string]bool)
    for _, issuance := range issuances {
        if dnsNames, ok := issuance["dns_names"].([]interface{}); ok {
            for _, name := range dnsNames {
                if nameStr, ok := name.(string); ok {
                    d := strings.TrimPrefix(nameStr, "*.")
                    if strings.Contains(d, domain) {
                        domains[d] = true
                    }
                }
            }
        }
    }
    
    result := make([]string, 0, len(domains))
    for d := range domains {
        result = append(result, d)
    }
    
    ctc.results[domain] = append(ctc.results[domain], result...)
    return result, nil
}

// CollectFromFacebookCT collects from Facebook CT
func (ctc *CertTransparencyCollector) CollectFromFacebookCT(domain string) ([]string, error) {
    url := fmt.Sprintf("https://ct.googleapis.com/icarus/ct/v1/get-entries?start=0&end=1000")
    // Facebook CT API is more complex, simplified for now
    return []string{}, nil
}

// CollectAll collects from all CT sources
func (ctc *CertTransparencyCollector) CollectAll(domain string) ([]string, error) {
    allDomains := make(map[string]bool)
    
    // crt.sh
    if domains, err := ctc.CollectFromCrtSh(domain); err == nil {
        for _, d := range domains {
            allDomains[d] = true
        }
    }
    
    // CertSpotter
    if domains, err := ctc.CollectFromCertSpotter(domain); err == nil {
        for _, d := range domains {
            allDomains[d] = true
        }
    }
    
    result := make([]string, 0, len(allDomains))
    for d := range allDomains {
        result = append(result, d)
    }
    
    return result, nil
}

// GetResults returns all collected domains
func (ctc *CertTransparencyCollector) GetResults() map[string][]string {
    return ctc.results
}

// Clear clears results
func (ctc *CertTransparencyCollector) Clear() {
    ctc.results = make(map[string][]string)
}
