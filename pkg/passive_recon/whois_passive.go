package passive_recon

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "strings"
    "time"
)

// WHOISInfo represents WHOIS information
type WHOISInfo struct {
    Domain      string    `json:"domain"`
    Registrar   string    `json:"registrar"`
    CreatedDate time.Time `json:"created_date"`
    ExpiryDate  time.Time `json:"expiry_date"`
    NameServers []string  `json:"name_servers"`
    Owner       string    `json:"owner"`
    Email       string    `json:"email"`
    Phone       string    `json:"phone"`
    Country     string    `json:"country"`
}

// PassiveWHOISCollector collects WHOIS data passively
type PassiveWHOISCollector struct {
    client  *http.Client
    results map[string]*WHOISInfo
}

// NewPassiveWHOISCollector creates a new WHOIS collector
func NewPassiveWHOISCollector() *PassiveWHOISCollector {
    return &PassiveWHOISCollector{
        client: &http.Client{Timeout: 30 * time.Second},
        results: make(map[string]*WHOISInfo),
    }
}

// CollectFromWhoisXMLAPI collects from WhoisXML API
func (pwc *PassiveWHOISCollector) CollectFromWhoisXMLAPI(domain, apiKey string) (*WHOISInfo, error) {
    url := fmt.Sprintf("https://www.whoisxmlapi.com/whoisserver/WhoisService?apiKey=%s&domainName=%s&outputFormat=JSON", apiKey, domain)
    resp, err := pwc.client.Get(url)
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
    
    whoisInfo := &WHOISInfo{
        Domain: domain,
    }
    
    // Parse response (simplified)
    if whoisResult, ok := result["WhoisRecord"].(map[string]interface{}); ok {
        if registrar, ok := whoisResult["registrarName"].(string); ok {
            whoisInfo.Registrar = registrar
        }
        if createdDate, ok := whoisResult["createdDate"].(string); ok {
            // Parse date
        }
    }
    
    pwc.results[domain] = whoisInfo
    return whoisInfo, nil
}

// CollectFromWhoisRDAP collects from RDAP
func (pwc *PassiveWHOISCollector) CollectFromWhoisRDAP(domain string) (*WHOISInfo, error) {
    url := fmt.Sprintf("https://rdap.org/domain/%s", domain)
    resp, err := pwc.client.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var result map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }
    
    whoisInfo := &WHOISInfo{
        Domain: domain,
    }
    
    // Parse RDAP response (simplified)
    if handle, ok := result["handle"].(string); ok {
        whoisInfo.Registrar = handle
    }
    
    pwc.results[domain] = whoisInfo
    return whoisInfo, nil
}

// CollectFromWhoisFreaks collects from WhoisFreaks API
func (pwc *PassiveWHOISCollector) CollectFromWhoisFreaks(domain, apiKey string) (*WHOISInfo, error) {
    url := fmt.Sprintf("https://api.whoisfreaks.com/v1.0/whois?apiKey=%s&domainName=%s", apiKey, domain)
    resp, err := pwc.client.Get(url)
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
    
    whoisInfo := &WHOISInfo{
        Domain: domain,
    }
    
    // Parse response
    if data, ok := result["data"].(map[string]interface{}); ok {
        if registrar, ok := data["registrar"].(string); ok {
            whoisInfo.Registrar = registrar
        }
    }
    
    pwc.results[domain] = whoisInfo
    return whoisInfo, nil
}

// GetResults returns all WHOIS results
func (pwc *PassiveWHOISCollector) GetResults() map[string]*WHOISInfo {
    return pwc.results
}

// ExtractEmails extracts emails from WHOIS data
func (pwc *PassiveWHOISCollector) ExtractEmails(domain string) []string {
    var emails []string
    if info, ok := pwc.results[domain]; ok && info.Email != "" {
        emails = append(emails, info.Email)
    }
    return emails
}

// ExtractNameServers extracts name servers
func (pwc *PassiveWHOISCollector) ExtractNameServers(domain string) []string {
    if info, ok := pwc.results[domain]; ok {
        return info.NameServers
    }
    return []string{}
}
