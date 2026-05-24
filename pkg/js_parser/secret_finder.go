package js_parser

import (
    "regexp"
    "strings"
)

// SecretType represents the type of secret found
type SecretType string

const (
    SecretJWT      SecretType = "JWT"
    SecretAWSKey   SecretType = "AWS_KEY"
    SecretAPIKey   SecretType = "API_KEY"
    SecretFirebase SecretType = "FIREBASE_CONFIG"
    SecretBearer   SecretType = "BEARER_TOKEN"
    SecretGeneric  SecretType = "GENERIC_SECRET"
)

// Secret represents a found secret
type Secret struct {
    Type    SecretType `json:"type"`
    Value   string     `json:"value"`
    Line    int        `json:"line"`
    Context string     `json:"context"`
}

// SecretFinder finds secrets in JavaScript
type SecretFinder struct {
    patterns map[SecretType]*regexp.Regexp
    secrets  []Secret
}

// NewSecretFinder creates a new secret finder
func NewSecretFinder() *SecretFinder {
    return &SecretFinder{
        patterns: map[SecretType]*regexp.Regexp{
            SecretJWT:      regexp.MustCompile(`eyJ[a-zA-Z0-9_-]{10,}\.[a-zA-Z0-9_-]{10,}\.[a-zA-Z0-9_-]{10,}`),
            SecretAWSKey:   regexp.MustCompile(`AKIA[0-9A-Z]{16}`),
            SecretAPIKey:   regexp.MustCompile(`[a-zA-Z0-9]{32,40}`),
            SecretFirebase: regexp.MustCompile(`[a-zA-Z0-9]+\.firebaseio\.com`),
            SecretBearer:   regexp.MustCompile(`Bearer\s+[a-zA-Z0-9_\-\.]+`),
            SecretGeneric:  regexp.MustCompile(`(?:api[_-]?key|secret|token|password|auth)\s*[:=]\s*['"]([^'"]{10,})['"]`),
        },
        secrets: make([]Secret, 0),
    }
}

// Find searches for secrets in source code
func (f *SecretFinder) Find(source string) []Secret {
    f.secrets = make([]Secret, 0)
    lines := strings.Split(source, "\n")
    
    for lineNum, line := range lines {
        for secretType, pattern := range f.patterns {
            matches := pattern.FindAllStringSubmatch(line, -1)
            for _, match := range matches {
                value := match[0]
                if len(match) > 1 {
                    value = match[1]
                }
                
                // Filter out false positives
                if f.isValidSecret(value) {
                    f.secrets = append(f.secrets, Secret{
                        Type:    secretType,
                        Value:   f.maskSecret(value),
                        Line:    lineNum + 1,
                        Context: strings.TrimSpace(line),
                    })
                }
            }
        }
    }
    
    return f.secrets
}

// isValidSecret checks if a secret is likely valid
func (f *SecretFinder) isValidSecret(value string) bool {
    // Ignore common false positives
    falsePositives := []string{
        "your-api-key", "your-secret", "example", "test", "demo",
        "api_key_here", "secret_here", "token_here", "xxxxx",
        "xxxxxxxx", "***", "SECRET", "API_KEY", "PASSWORD",
    }
    
    lowerValue := strings.ToLower(value)
    for _, fp := range falsePositives {
        if strings.Contains(lowerValue, fp) {
            return false
        }
    }
    
    // Must have at least 10 characters
    if len(value) < 10 {
        return false
    }
    
    return true
}

// maskSecret masks the secret for safe display
func (f *SecretFinder) maskSecret(secret string) string {
    if len(secret) <= 8 {
        return "***"
    }
    return secret[:4] + "..." + secret[len(secret)-4:]
}

// GetSecrets returns all found secrets
func (f *SecretFinder) GetSecrets() []Secret {
    return f.secrets
}

// GetSecretsByType returns secrets of a specific type
func (f *SecretFinder) GetSecretsByType(secretType SecretType) []Secret {
    var filtered []Secret
    for _, s := range f.secrets {
        if s.Type == secretType {
            filtered = append(filtered, s)
        }
    }
    return filtered
}
