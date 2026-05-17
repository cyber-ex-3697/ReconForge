package vulnerability

import (
    "regexp"
    "strings"
)

type FalsePositiveDetector struct {
    patterns []string
    whitelist []string
}

func NewFalsePositiveDetector() *FalsePositiveDetector {
    return &FalsePositiveDetector{
        patterns: []string{
            `example\.com`,
            `test\.com`,
            `localhost`,
            `127\.0\.0\.1`,
            `0\.0\.0\.0`,
            `demo\.`,
            `sample\.`,
            `placeholder`,
            `TODO`,
            `FIXME`,
        },
        whitelist: []string{
            `google\.com`,
            `cloudflare\.com`,
            `akamai\.net`,
        },
    }
}

func (f *FalsePositiveDetector) IsFalsePositive(vuln Vulnerability) bool {
    // Check patterns
    for _, pattern := range f.patterns {
        re := regexp.MustCompile(pattern)
        if re.MatchString(vuln.URL) {
            return true
        }
        if re.MatchString(vuln.Name) {
            return true
        }
    }
    
    // Check if whitelisted (not false positive)
    for _, wl := range f.whitelist {
        re := regexp.MustCompile(wl)
        if re.MatchString(vuln.URL) {
            return false
        }
    }
    
    // Check CVSS score - very low score might be false positive
    if vuln.CVSS < 2.0 {
        return true
    }
    
    return false
}

func (f *FalsePositiveDetector) Filter(vulns []Vulnerability) []Vulnerability {
    var filtered []Vulnerability
    for _, v := range vulns {
        if !f.IsFalsePositive(v) {
            filtered = append(filtered, v)
        }
    }
    return filtered
}

func (f *FalsePositiveDetector) AddPattern(pattern string) {
    f.patterns = append(f.patterns, pattern)
}

func (f *FalsePositiveDetector) AddWhitelist(pattern string) {
    f.whitelist = append(f.whitelist, pattern)
}

func (f *FalsePositiveDetector) IsKnownFalsePositive(templateID string) bool {
    knownFalsePositives := []string{
        "http-missing-security-headers",
        "http-cors-misconfig",
        "http-missing-security-tag",
    }
    
    for _, fp := range knownFalsePositives {
        if strings.Contains(templateID, fp) {
            return true
        }
    }
    return false
}
