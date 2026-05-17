package takeover

import (
    "net"
    "strings"
)

type CNAMEChecker struct {
    vulnerableServices []string
}

func NewCNAMEChecker() *CNAMEChecker {
    return &CNAMEChecker{
        vulnerableServices: []string{
            "amazonaws.com",
            "azurewebsites.net",
            "cloudfront.net",
            "herokuapp.com",
            "github.io",
            "readme.io",
            "surge.sh",
            "bitbucket.io",
        },
    }
}

func (c *CNAMEChecker) Check(domain string) (bool, string, string) {
    cname, err := net.LookupCNAME(domain)
    if err != nil {
        return false, "", ""
    }
    
    cname = strings.TrimSuffix(cname, ".")
    
    for _, service := range c.vulnerableServices {
        if strings.Contains(cname, service) {
            return true, cname, service
        }
    }
    return false, cname, ""
}

func (c *CNAMEChecker) CheckBatch(domains []string) map[string]string {
    results := make(map[string]string)
    for _, domain := range domains {
        vulnerable, cname, service := c.Check(domain)
        if vulnerable {
            results[domain] = service + " -> " + cname
        }
    }
    return results
}
