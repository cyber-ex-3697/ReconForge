package auth_testing

import (
    "fmt"
    "io"
    "net/http"
    "net/url"
    "regexp"
    "strings"
    "time"
)

// OAuthChecker checks OAuth misconfigurations
type OAuthChecker struct {
    client     *http.Client
    findings   []OAuthFinding
}

// OAuthFinding represents an OAuth misconfiguration finding
type OAuthFinding struct {
    Type        string `json:"type"`
    URL         string `json:"url"`
    Parameter   string `json:"parameter"`
    Issue       string `json:"issue"`
    Severity    string `json:"severity"`
}

// NewOAuthChecker creates a new OAuth checker
func NewOAuthChecker() *OAuthChecker {
    return &OAuthChecker{
        client: &http.Client{
            Timeout: 15 * time.Second,
            CheckRedirect: func(req *http.Request, via []*http.Request) error {
                return nil
            },
        },
        findings: make([]OAuthFinding, 0),
    }
}

// CheckOAuthEndpoint checks an OAuth endpoint for misconfigurations
func (oc *OAuthChecker) CheckOAuthEndpoint(endpointURL string) ([]OAuthFinding, error) {
    oc.findings = make([]OAuthFinding, 0)
    
    // Check for redirect URI manipulation
    oc.checkRedirectURIManipulation(endpointURL)
    
    // Check for response type manipulation
    oc.checkResponseTypeManipulation(endpointURL)
    
    // Check for scope escalation
    oc.checkScopeEscalation(endpointURL)
    
    // Check for client ID exposure
    oc.checkClientIDExposure(endpointURL)
    
    return oc.findings, nil
}

// checkRedirectURIManipulation tests redirect URI validation
func (oc *OAuthChecker) checkRedirectURIManipulation(endpointURL string) {
    testCases := []string{
        "https://evil.com",
        "https://attacker.com/callback",
        "https://example.com.evil.com",
        "https://example.com@evil.com",
        "javascript://example.com/%0Aalert(1)",
        "data:text/html,<script>alert(1)</script>",
    }
    
    baseURL, err := url.Parse(endpointURL)
    if err != nil {
        return
    }
    
    for _, testURI := range testCases {
        params := baseURL.Query()
        params.Set("redirect_uri", testURI)
        baseURL.RawQuery = params.Encode()
        
        resp, err := oc.client.Get(baseURL.String())
        if err != nil {
            continue
        }
        defer resp.Body.Close()
        
        if resp.StatusCode == 302 || resp.StatusCode == 301 {
            location := resp.Header.Get("Location")
            if strings.Contains(location, testURI) {
                oc.findings = append(oc.findings, OAuthFinding{
                    Type:      "redirect_uri",
                    URL:       endpointURL,
                    Parameter: "redirect_uri",
                    Issue:     fmt.Sprintf("Open redirect to: %s", testURI),
                    Severity:  "HIGH",
                })
            }
        }
        
        time.Sleep(100 * time.Millisecond)
    }
}

// checkResponseTypeManipulation tests response type validation
func (oc *OAuthChecker) checkResponseTypeManipulation(endpointURL string) {
    testTypes := []string{
        "token",
        "code token",
        "id_token",
        "id_token token",
        "code id_token",
        "code id_token token",
    }
    
    baseURL, err := url.Parse(endpointURL)
    if err != nil {
        return
    }
    
    for _, respType := range testTypes {
        params := baseURL.Query()
        params.Set("response_type", respType)
        baseURL.RawQuery = params.Encode()
        
        resp, err := oc.client.Get(baseURL.String())
        if err != nil {
            continue
        }
        defer resp.Body.Close()
        
        if resp.StatusCode == 200 {
            oc.findings = append(oc.findings, OAuthFinding{
                Type:      "response_type",
                URL:       endpointURL,
                Parameter: "response_type",
                Issue:     fmt.Sprintf("Response type '%s' accepted", respType),
                Severity:  "MEDIUM",
            })
        }
        
        time.Sleep(100 * time.Millisecond)
    }
}

// checkScopeEscalation tests scope escalation
func (oc *OAuthChecker) checkScopeEscalation(endpointURL string) {
    testScopes := []string{
        "admin",
        "admin:*",
        "user:write",
        "delete",
        "modify",
        "full_access",
        "all",
        "*",
    }
    
    baseURL, err := url.Parse(endpointURL)
    if err != nil {
        return
    }
    
    for _, scope := range testScopes {
        params := baseURL.Query()
        params.Set("scope", scope)
        baseURL.RawQuery = params.Encode()
        
        resp, err := oc.client.Get(baseURL.String())
        if err != nil {
            continue
        }
        defer resp.Body.Close()
        
        if resp.StatusCode == 200 {
            oc.findings = append(oc.findings, OAuthFinding{
                Type:      "scope",
                URL:       endpointURL,
                Parameter: "scope",
                Issue:     fmt.Sprintf("Scope '%s' accepted without validation", scope),
                Severity:  "CRITICAL",
            })
        }
        
        time.Sleep(100 * time.Millisecond)
    }
}

// checkClientIDExposure checks for client ID exposure
func (oc *OAuthChecker) checkClientIDExposure(endpointURL string) {
    baseURL, err := url.Parse(endpointURL)
    if err != nil {
        return
    }
    
    params := baseURL.Query()
    params.Del("client_id")
    baseURL.RawQuery = params.Encode()
    
    resp, err := oc.client.Get(baseURL.String())
    if err != nil {
        return
    }
    defer resp.Body.Close()
    
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return
    }
    
    // Look for client ID patterns in response
    clientIDPattern := regexp.MustCompile(`[a-zA-Z0-9]{20,}`)
    matches := clientIDPattern.FindAllString(string(body), -1)
    
    for _, match := range matches {
        oc.findings = append(oc.findings, OAuthFinding{
            Type:      "client_id_exposure",
            URL:       endpointURL,
            Parameter: "response_body",
            Issue:     fmt.Sprintf("Potential client ID exposed: %s", match),
            Severity:  "MEDIUM",
        })
    }
}

// GetFindings returns all OAuth findings
func (oc *OAuthChecker) GetFindings() []OAuthFinding {
    return oc.findings
}
