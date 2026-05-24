package auth_testing

import (
    "bytes"
    "crypto/md5"
    "encoding/hex"
    "fmt"
    "io"
    "net/http"
    "regexp"
    "strconv"
    "strings"
    "time"
)

// IDORDetector detects Insecure Direct Object Reference vulnerabilities
type IDORDetector struct {
    client       *http.Client
    patterns     []*regexp.Regexp
    results      []IDORFinding
    session1     *Session
    session2     *Session
}

// IDORFinding represents an IDOR vulnerability finding
type IDORFinding struct {
    URL          string            `json:"url"`
    Parameter    string            `json:"parameter"`
    OriginalValue string           `json:"original_value"`
    ModifiedValue string           `json:"modified_value"`
    OriginalResponse string        `json:"original_response_hash"`
    ModifiedResponse string        `json:"modified_response_hash"`
    StatusCode   int               `json:"status_code"`
    Confidence   float64           `json:"confidence"`
    Exploitable  bool              `json:"exploitable"`
}

// Session represents a user session
type Session struct {
    ID       string
    Cookies  []*http.Cookie
    Headers  map[string]string
    UserRole string
}

// NewIDORDetector creates a new IDOR detector
func NewIDORDetector() *IDORDetector {
    return &IDORDetector{
        client: &http.Client{
            Timeout: 30 * time.Second,
            CheckRedirect: func(req *http.Request, via []*http.Request) error {
                return nil // Don't follow redirects
            },
        },
        patterns: []*regexp.Regexp{
            regexp.MustCompile(`[?&](id|user_id|account_id|profile_id|doc_id|file_id|order_id|invoice_id)=(\d+)`),
            regexp.MustCompile(`[?&](userId|accountId|profileId|documentId)=([a-f0-9-]+)`),
            regexp.MustCompile(`/api/(?:users|accounts|profiles|documents)/(\d+)`),
            regexp.MustCompile(`/api/(?:users|accounts|profiles|documents)/([a-f0-9-]+)`),
        },
        results: make([]IDORFinding, 0),
    }
}

// SetSessions sets two different user sessions for testing
func (d *IDORDetector) SetSessions(session1, session2 *Session) {
    d.session1 = session1
    d.session2 = session2
}

// ScanURL scans a URL for IDOR vulnerabilities
func (d *IDORDetector) ScanURL(url string) ([]IDORFinding, error) {
    d.results = make([]IDORFinding, 0)
    
    // Extract parameters from URL
    for _, pattern := range d.patterns {
        matches := pattern.FindAllStringSubmatch(url, -1)
        for _, match := range matches {
            if len(match) >= 3 {
                finding := d.testParameter(url, match[1], match[2])
                if finding.Confidence > 0.5 {
                    d.results = append(d.results, finding)
                }
            }
        }
    }
    
    return d.results, nil
}

// testParameter tests a single parameter for IDOR
func (d *IDORDetector) testParameter(url, param, value string) IDORFinding {
    finding := IDORFinding{
        URL:           url,
        Parameter:     param,
        OriginalValue: value,
        Confidence:    0.0,
        Exploitable:   false,
    }
    
    // Generate variations of the value
    variations := d.generateVariations(value)
    
    for _, variant := range variations {
        // Modify URL with new value
        modifiedURL := strings.Replace(url, param+"="+value, param+"="+variant, 1)
        
        // Test with first session
        resp1, body1 := d.makeRequest(modifiedURL, d.session1)
        if resp1 == nil {
            continue
        }
        
        // Test with second session (if available)
        if d.session2 != nil {
            resp2, body2 := d.makeRequest(modifiedURL, d.session2)
            if resp2 != nil {
                // Compare responses
                similarity := d.compareResponses(body1, body2)
                
                // Calculate confidence
                confidence := d.calculateConfidence(resp1, resp2, similarity)
                if confidence > finding.Confidence {
                    finding.Confidence = confidence
                    finding.ModifiedValue = variant
                    finding.OriginalResponseHash = d.hashResponse(body1)
                    finding.ModifiedResponseHash = d.hashResponse(body2)
                    finding.StatusCode = resp2.StatusCode
                    finding.Exploitable = confidence > 0.7 && resp2.StatusCode == 200
                }
            }
        }
    }
    
    return finding
}

// generateVariations generates value variations for IDOR testing
func (d *IDORDetector) generateVariations(original string) []string {
    variations := make([]string, 0)
    
    // Try numeric increments/decrements
    if num, err := strconv.Atoi(original); err == nil {
        variations = append(variations, strconv.Itoa(num+1))
        variations = append(variations, strconv.Itoa(num-1))
        variations = append(variations, strconv.Itoa(num+10))
        variations = append(variations, strconv.Itoa(num-10))
        variations = append(variations, strconv.Itoa(num+100))
        variations = append(variations, strconv.Itoa(num-100))
        variations = append(variations, "1")
        variations = append(variations, "0")
        variations = append(variations, "999999")
    }
    
    // Try UUID variations
    if strings.Contains(original, "-") {
        parts := strings.Split(original, "-")
        if len(parts) == 5 {
            // Modify last segment
            if last, err := strconv.ParseInt(parts[4], 16, 64); err == nil {
                newLast := fmt.Sprintf("%x", last+1)
                variations = append(variations, strings.Join(append(parts[:4], newLast), "-"))
            }
            // Use all zeros
            variations = append(variations, "00000000-0000-0000-0000-000000000000")
            // Use all ones
            variations = append(variations, "11111111-1111-1111-1111-111111111111")
            // Use f's
            variations = append(variations, "ffffffff-ffff-ffff-ffff-ffffffffffff")
        }
    }
    
    // Try empty/null values
    variations = append(variations, "")
    variations = append(variations, "null")
    variations = append(variations, "NULL")
    variations = append(variations, "0")
    variations = append(variations, "-1")
    
    return variations
}

// makeRequest makes an HTTP request with a session
func (d *IDORDetector) makeRequest(url string, session *Session) (*http.Response, string) {
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, ""
    }
    
    // Add session cookies
    for _, cookie := range session.Cookies {
        req.AddCookie(cookie)
    }
    
    // Add session headers
    for k, v := range session.Headers {
        req.Header.Set(k, v)
    }
    
    resp, err := d.client.Do(req)
    if err != nil {
        return nil, ""
    }
    defer resp.Body.Close()
    
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return resp, ""
    }
    
    return resp, string(body)
}

// compareResponses compares two responses using MinHash-like approach
func (d *IDORDetector) compareResponses(body1, body2 string) float64 {
    if body1 == body2 {
        return 1.0
    }
    
    // Simple similarity based on length and content
    len1 := len(body1)
    len2 := len(body2)
    
    if len1 == 0 && len2 == 0 {
        return 1.0
    }
    
    // Length ratio
    lenRatio := 1.0 - float64(abs(len1-len2))/float64(max(len1, len2))
    
    // Word set similarity
    words1 := strings.Fields(body1)
    words2 := strings.Fields(body2)
    
    set1 := make(map[string]bool)
    set2 := make(map[string]bool)
    
    for _, w := range words1 {
        set1[w] = true
    }
    for _, w := range words2 {
        set2[w] = true
    }
    
    intersection := 0
    for w := range set1 {
        if set2[w] {
            intersection++
        }
    }
    
    union := len(set1) + len(set2) - intersection
    wordSimilarity := 0.0
    if union > 0 {
        wordSimilarity = float64(intersection) / float64(union)
    }
    
    // Combine similarities
    return (lenRatio + wordSimilarity) / 2.0
}

// calculateConfidence calculates confidence score for IDOR finding
func (d *IDORDetector) calculateConfidence(resp1, resp2 *http.Response, similarity float64) float64 {
    confidence := 0.0
    
    // Status code matters
    if resp2.StatusCode == 200 {
        confidence += 0.3
    } else if resp2.StatusCode == 403 || resp2.StatusCode == 401 {
        confidence -= 0.2
    }
    
    // Different responses indicate different data
    if similarity < 0.3 {
        confidence += 0.4
    } else if similarity < 0.6 {
        confidence += 0.2
    }
    
    // Content length difference
    if resp1.ContentLength != resp2.ContentLength && resp2.ContentLength > 0 {
        confidence += 0.2
    }
    
    return confidence
}

// hashResponse creates a hash of response body
func (d *IDORDetector) hashResponse(body string) string {
    hash := md5.Sum([]byte(body))
    return hex.EncodeToString(hash[:])
}

func abs(x int) int {
    if x < 0 {
        return -x
    }
    return x
}

func max(a, b int) int {
    if a > b {
        return a
    }
    return b
}
