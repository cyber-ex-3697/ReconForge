package auth_testing

import (
    "encoding/base64"
    "encoding/json"
    "fmt"
    "strings"
    "time"
)

// JWTAnalyzer analyzes JWT tokens
type JWTAnalyzer struct {
    token       string
    header      map[string]interface{}
    payload     map[string]interface{}
    signature   string
    isValid     bool
}

// JWTFinding represents a JWT vulnerability finding
type JWTFinding struct {
    Issue       string   `json:"issue"`
    Severity    string   `json:"severity"`
    Description string   `json:"description"`
    Recommendation string `json:"recommendation"`
}

// NewJWTAnalyzer creates a new JWT analyzer
func NewJWTAnalyzer(token string) *JWTAnalyzer {
    ja := &JWTAnalyzer{
        token:   token,
        header:  make(map[string]interface{}),
        payload: make(map[string]interface{}),
    }
    ja.parse()
    return ja
}

// parse parses the JWT token
func (ja *JWTAnalyzer) parse() {
    parts := strings.Split(ja.token, ".")
    if len(parts) != 3 {
        ja.isValid = false
        return
    }
    
    // Decode header
    headerJSON, err := base64.RawURLEncoding.DecodeString(parts[0])
    if err == nil {
        json.Unmarshal(headerJSON, &ja.header)
    }
    
    // Decode payload
    payloadJSON, err := base64.RawURLEncoding.DecodeString(parts[1])
    if err == nil {
        json.Unmarshal(payloadJSON, &ja.payload)
    }
    
    ja.signature = parts[2]
    ja.isValid = true
}

// Analyze performs comprehensive JWT analysis
func (ja *JWTAnalyzer) Analyze() []JWTFinding {
    findings := make([]JWTFinding, 0)
    
    // Check for none algorithm
    if alg, ok := ja.header["alg"].(string); ok {
        if alg == "none" {
            findings = append(findings, JWTFinding{
                Issue:       "None Algorithm",
                Severity:    "CRITICAL",
                Description: "JWT uses 'none' algorithm which allows arbitrary token forgery",
                Recommendation: "Remove 'none' algorithm support. Use strong algorithms like HS256 or RS256.",
            })
        }
        
        // Check for weak algorithm
        if alg == "HS256" {
            findings = append(findings, JWTFinding{
                Issue:       "Weak Algorithm",
                Severity:    "MEDIUM",
                Description: "JWT uses HMAC which may be vulnerable to key brute-forcing",
                Recommendation: "Consider using RS256 with strong key management.",
            })
        }
    }
    
    // Check expiration
    if exp, ok := ja.payload["exp"].(float64); ok {
        expTime := time.Unix(int64(exp), 0)
        if expTime.Before(time.Now()) {
            findings = append(findings, JWTFinding{
                Issue:       "Expired Token",
                Severity:    "LOW",
                Description: "Token has expired",
                Recommendation: "Token has expired and should be refreshed.",
            })
        } else if expTime.Before(time.Now().Add(24 * time.Hour)) {
            findings = append(findings, JWTFinding{
                Issue:       "Short Expiration",
                Severity:    "INFO",
                Description: "Token expires soon",
                Recommendation: "Token expiration is reasonable for security.",
            })
        }
    } else {
        findings = append(findings, JWTFinding{
            Issue:       "No Expiration",
            Severity:    "HIGH",
            Description: "JWT does not have an expiration claim",
            Recommendation: "Add 'exp' claim to enforce token expiration.",
        })
    }
    
    // Check for sensitive data in payload
    sensitiveFields := []string{"password", "secret", "key", "token", "credit", "ssn"}
    for field := range ja.payload {
        for _, sensitive := range sensitiveFields {
            if strings.Contains(strings.ToLower(field), sensitive) {
                findings = append(findings, JWTFinding{
                    Issue:       "Sensitive Data Exposure",
                    Severity:    "HIGH",
                    Description: fmt.Sprintf("Sensitive field '%s' found in JWT payload", field),
                    Recommendation: "Do not store sensitive data in JWT payload.",
                })
                break
            }
        }
    }
    
    return findings
}

// GetHeader returns JWT header
func (ja *JWTAnalyzer) GetHeader() map[string]interface{} {
    return ja.header
}

// GetPayload returns JWT payload
func (ja *JWTAnalyzer) GetPayload() map[string]interface{} {
    return ja.payload
}

// GetSignature returns JWT signature
func (ja *JWTAnalyzer) GetSignature() string {
    return ja.signature
}

// IsValid returns whether token is validly formatted
func (ja *JWTAnalyzer) IsValid() bool {
    return ja.isValid
}

// GetKidClaim returns kid claim if present
func (ja *JWTAnalyzer) GetKidClaim() string {
    if kid, ok := ja.header["kid"].(string); ok {
        return kid
    }
    return ""
}

// GetIssuer returns issuer claim
func (ja *JWTAnalyzer) GetIssuer() string {
    if iss, ok := ja.payload["iss"].(string); ok {
        return iss
    }
    return ""
}

// GetSubject returns subject claim
func (ja *JWTAnalyzer) GetSubject() string {
    if sub, ok := ja.payload["sub"].(string); ok {
        return sub
    }
    return ""
}
