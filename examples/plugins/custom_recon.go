package main

import (
    "bufio"
    "context"
    "fmt"
    "os"
    "os/exec"
    "strings"
    
    "reconforge/plugins"
)

// =============================================================================
// CustomReconPlugin - Example Custom Reconnaissance Plugin
// =============================================================================

type CustomReconPlugin struct {
    plugins.BasePlugin
    config map[string]interface{}
}

// CustomResult represents the plugin's output
type CustomResult struct {
    CustomSubdomains []string `json:"custom_subdomains"`
    CustomURLs       []string `json:"custom_urls"`
    CustomFindings   []string `json:"custom_findings"`
}

// NewCustomReconPlugin creates a new instance
func NewCustomReconPlugin() *CustomReconPlugin {
    return &CustomReconPlugin{
        BasePlugin: plugins.BasePlugin{
            Info: plugins.PluginInfo{
                Name:        "CustomReconPlugin",
                Version:     "1.0.0",
                Author:      "ReconForge Security Team",
                Description: "Custom reconnaissance using additional techniques",
                Phase:       "subdomain",
                Dependencies: []string{"curl", "jq"},
            },
        },
    }
}

// Init initializes the plugin with configuration
func (p *CustomReconPlugin) Init(config map[string]interface{}) error {
    p.config = config
    return nil
}

// Run executes the custom reconnaissance
func (p *CustomReconPlugin) Run(ctx context.Context, target string, input interface{}) (*plugins.PluginResult, error) {
    fmt.Printf("[Plugin] %s v%s starting...\n", p.Info.Name, p.Info.Version)
    fmt.Printf("[Plugin] Target: %s\n", target)
    
    result := &CustomResult{
        CustomSubdomains: make([]string, 0),
        CustomURLs:       make([]string, 0),
        CustomFindings:   make([]string, 0),
    }
    
    // =========================================================================
    // 1. Custom Subdomain Discovery
    // =========================================================================
    fmt.Println("[Plugin] Discovering custom subdomains...")
    
    // Common subdomain prefixes for manual testing
    prefixes := []string{
        "admin", "api", "app", "auth", "backup", "beta", "blog", "cdn",
        "dashboard", "dev", "docs", "email", "ftp", "git", "internal",
        "login", "mail", "old", "portal", "remote", "secure", "stage",
        "staging", "support", "test", "vpn", "webmail", "www",
    }
    
    for _, prefix := range prefixes {
        subdomain := fmt.Sprintf("%s.%s", prefix, target)
        if p.checkSubdomain(subdomain) {
            result.CustomSubdomains = append(result.CustomSubdomains, subdomain)
            fmt.Printf("  Found: %s\n", subdomain)
        }
    }
    
    // =========================================================================
    // 2. Custom URL Discovery
    // =========================================================================
    fmt.Println("[Plugin] Discovering custom URLs...")
    
    // Common sensitive paths
    paths := []string{
        "/admin", "/backup", "/config", "/console", "/dashboard",
        "/debug", "/dev", "/docs", "/dump", "/logs", "/phpinfo",
        "/private", "/robots.txt", "/sitemap.xml", "/swagger",
        "/test", "/tmp", "/wp-admin", "/wp-config.php.bak",
    }
    
    baseURL := fmt.Sprintf("https://%s", target)
    for _, path := range paths {
        url := baseURL + path
        if p.checkURL(url) {
            result.CustomURLs = append(result.CustomURLs, url)
            fmt.Printf("  Found URL: %s\n", url)
        }
    }
    
    // =========================================================================
    // 3. Custom Security Checks
    // =========================================================================
    fmt.Println("[Plugin] Running custom security checks...")
    
    // Check for exposed git repository
    gitURL := baseURL + "/.git/config"
    if p.checkURL(gitURL) {
        result.CustomFindings = append(result.CustomFindings, 
            fmt.Sprintf("Exposed .git repository: %s", gitURL))
        fmt.Printf("  [!] Exposed .git repository found!\n")
    }
    
    // Check for exposed environment file
    envURL := baseURL + "/.env"
    if p.checkURL(envURL) {
        result.CustomFindings = append(result.CustomFindings,
            fmt.Sprintf("Exposed .env file: %s", envURL))
        fmt.Printf("  [!] Exposed .env file found!\n")
    }
    
    // Check for exposed backup files
    backupPatterns := []string{".bak", ".backup", ".old", ".sql", ".tar.gz", ".zip"}
    for _, pattern := range backupPatterns {
        backupURL := baseURL + "/backup" + pattern
        if p.checkURL(backupURL) {
            result.CustomFindings = append(result.CustomFindings,
                fmt.Sprintf("Backup file: %s", backupURL))
            fmt.Printf("  [!] Backup file found: %s\n", pattern)
            break
        }
    }
    
    // =========================================================================
    // 4. Save Results
    // =========================================================================
    outputFile := p.OutputDir + "/custom_recon_results.txt"
    content := p.formatResults(result, target)
    os.WriteFile(outputFile, []byte(content), 0644)
    
    fmt.Printf("[Plugin] Results saved to: %s\n", outputFile)
    
    return &plugins.PluginResult{
        Success:    true,
        Data:       result,
        OutputFile: outputFile,
    }, nil
}

// checkSubdomain checks if a subdomain resolves
func (p *CustomReconPlugin) checkSubdomain(subdomain string) bool {
    cmd := exec.Command("dig", "+short", subdomain)
    output, err := cmd.Output()
    if err != nil {
        return false
    }
    return len(strings.TrimSpace(string(output))) > 0
}

// checkURL checks if a URL is accessible
func (p *CustomReconPlugin) checkURL(url string) bool {
    cmd := exec.Command("curl", "-s", "-o", "/dev/null", "-w", "%{http_code}", "--max-time", "5", url)
    output, err := cmd.Output()
    if err != nil {
        return false
    }
    code := strings.TrimSpace(string(output))
    return code == "200" || code == "301" || code == "302" || code == "403"
}

// formatResults formats the results for output
func (p *CustomReconPlugin) formatResults(result *CustomResult, target string) string {
    var sb strings.Builder
    
    sb.WriteString("=" + strings.Repeat("=", 70) + "\n")
    sb.WriteString("  RECONFORGE - CUSTOM RECONNAISSANCE RESULTS\n")
    sb.WriteString("=" + strings.Repeat("=", 70) + "\n\n")
    
    sb.WriteString(fmt.Sprintf("Target: %s\n", target))
    sb.WriteString(fmt.Sprintf("Plugin: %s v%s\n", p.Info.Name, p.Info.Version))
    sb.WriteString(fmt.Sprintf("Time: %s\n\n", strings.Repeat("=", 50)))
    
    // Custom Subdomains
    sb.WriteString("\n[+] CUSTOM SUBDOMAINS FOUND\n")
    sb.WriteString(strings.Repeat("-", 50) + "\n")
    if len(result.CustomSubdomains) > 0 {
        for _, sub := range result.CustomSubdomains {
            sb.WriteString(fmt.Sprintf("  - %s\n", sub))
        }
    } else {
        sb.WriteString("  None found\n")
    }
    
    // Custom URLs
    sb.WriteString("\n[+] CUSTOM URLS FOUND\n")
    sb.WriteString(strings.Repeat("-", 50) + "\n")
    if len(result.CustomURLs) > 0 {
        for _, url := range result.CustomURLs {
            sb.WriteString(fmt.Sprintf("  - %s\n", url))
        }
    } else {
        sb.WriteString("  None found\n")
    }
    
    // Custom Findings
    sb.WriteString("\n[!] SECURITY FINDINGS\n")
    sb.WriteString(strings.Repeat("-", 50) + "\n")
    if len(result.CustomFindings) > 0 {
        for _, finding := range result.CustomFindings {
            sb.WriteString(fmt.Sprintf("  - %s\n", finding))
        }
    } else {
        sb.WriteString("  None found\n")
    }
    
    sb.WriteString("\n" + strings.Repeat("=", 70) + "\n")
    sb.WriteString("  Scan completed. Manual verification required for findings.\n")
    sb.WriteString(strings.Repeat("=", 70) + "\n")
    
    return sb.String()
}

// Cleanup performs cleanup operations
func (p *CustomReconPlugin) Cleanup() error {
    // Clean up any temporary files
    return nil
}

// =============================================================================
// Plugin Entry Point
// =============================================================================

// Plugin instance that will be loaded
var Plugin = NewCustomReconPlugin()

// Export for plugin system
func init() {
    // Plugin registration would happen here
}
