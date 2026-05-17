package utils

import (
    "regexp"
    "strings"
)

func SanitizeDomain(domain string) string {
    domain = strings.TrimPrefix(domain, "https://")
    domain = strings.TrimPrefix(domain, "http://")
    domain = strings.TrimPrefix(domain, "www.")
    domain = strings.Split(domain, "/")[0]
    domain = strings.Split(domain, ":")[0]
    
    re := regexp.MustCompile(`[^a-zA-Z0-9.-]`)
    domain = re.ReplaceAllString(domain, "")
    
    return strings.ToLower(domain)
}

func SanitizeFilename(filename string) string {
    re := regexp.MustCompile(`[^a-zA-Z0-9._-]`)
    return re.ReplaceAllString(filename, "_")
}

func IsValidDomain(domain string) bool {
    re := regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
    return re.MatchString(domain)
}
