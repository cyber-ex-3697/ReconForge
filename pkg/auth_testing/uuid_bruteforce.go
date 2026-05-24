package auth_testing

import (
    "fmt"
    "io"
    "net/http"
    "regexp"
    "strconv"
    "strings"
    "time"
)

// UUIDBruteforce performs UUID bruteforce heuristics
type UUIDBruteforce struct {
    client      *http.Client
    baseURL     string
    pattern     *regexp.Regexp
    results     []UUIDFinding
}

// UUIDFinding represents a UUID bruteforce finding
type UUIDFinding struct {
    URL         string   `json:"url"`
    UUIDPattern string   `json:"uuid_pattern"`
    ValidUUIDs  []string `json:"valid_uuids"`
    TotalTested int      `json:"total_tested"`
    SuccessRate float64  `json:"success_rate"`
}

// NewUUIDBruteforce creates a new UUID bruteforce handler
func NewUUIDBruteforce() *UUIDBruteforce {
    return &UUIDBruteforce{
        client: &http.Client{
            Timeout: 10 * time.Second,
        },
        pattern: regexp.MustCompile(`[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}`),
        results: make([]UUIDFinding, 0),
    }
}

// SetBaseURL sets the base URL for testing
func (ub *UUIDBruteforce) SetBaseURL(baseURL string) {
    ub.baseURL = strings.TrimSuffix(baseURL, "/")
}

// ExtractUUIDs extracts UUIDs from URL
func (ub *UUIDBruteforce) ExtractUUIDs(url string) []string {
    uuids := make([]string, 0)
    matches := ub.pattern.FindAllString(url, -1)
    
    for _, match := range matches {
        uuids = append(uuids, match)
    }
    
    return uuids
}

// BruteforceUUID tests UUID variations
func (ub *UUIDBruteforce) BruteforceUUID(url string, originalUUID string) (*UUIDFinding, error) {
    finding := &UUIDFinding{
        URL:         url,
        UUIDPattern: originalUUID,
        ValidUUIDs:  make([]string, 0),
    }
    
    variations := ub.generateUUIDVariations(originalUUID)
    finding.TotalTested = len(variations)
    
    for _, variant := range variations {
        modifiedURL := strings.Replace(url, originalUUID, variant, -1)
        
        resp, err := ub.client.Get(modifiedURL)
        if err != nil {
            continue
        }
        defer resp.Body.Close()
        
        if resp.StatusCode == 200 {
            finding.ValidUUIDs = append(finding.ValidUUIDs, variant)
        }
        
        time.Sleep(100 * time.Millisecond) // Rate limiting
    }
    
    if len(finding.ValidUUIDs) > 0 {
        finding.SuccessRate = float64(len(finding.ValidUUIDs)) / float64(finding.TotalTested)
        ub.results = append(ub.results, *finding)
    }
    
    return finding, nil
}

// generateUUIDVariations generates UUID variations for testing
func (ub *UUIDBruteforce) generateUUIDVariations(originalUUID string) []string {
    variations := make([]string, 0)
    parts := strings.Split(originalUUID, "-")
    
    if len(parts) != 5 {
        return variations
    }
    
    // Modify each segment
    for i, part := range parts {
        if num, err := strconv.ParseInt(part, 16, 64); err == nil {
            for delta := -5; delta <= 5; delta++ {
                newNum := num + int64(delta)
                if newNum >= 0 {
                    newPart := fmt.Sprintf("%x", newNum)
                    // Pad to correct length
                    for len(newPart) < len(part) {
                        newPart = "0" + newPart
                    }
                    if len(newPart) > len(part) {
                        newPart = newPart[:len(part)]
                    }
                    
                    newParts := make([]string, len(parts))
                    copy(newParts, parts)
                    newParts[i] = newPart
                    variations = append(variations, strings.Join(newParts, "-"))
                }
            }
        }
    }
    
    // Add common UUID variations
    variations = append(variations,
        "00000000-0000-0000-0000-000000000000",
        "11111111-1111-1111-1111-111111111111",
        "ffffffff-ffff-ffff-ffff-ffffffffffff",
        "99999999-9999-9999-9999-999999999999",
    )
    
    return variations
}

// ScanURL scans a URL for UUID-based IDOR
func (ub *UUIDBruteforce) ScanURL(url string) ([]UUIDFinding, error) {
    uuids := ub.ExtractUUIDs(url)
    results := make([]UUIDFinding, 0)
    
    for _, uuid := range uuids {
        finding, err := ub.BruteforceUUID(url, uuid)
        if err == nil && finding.SuccessRate > 0 {
            results = append(results, *finding)
        }
    }
    
    return results, nil
}

// GetResults returns all findings
func (ub *UUIDBruteforce) GetResults() []UUIDFinding {
    return ub.results
}
