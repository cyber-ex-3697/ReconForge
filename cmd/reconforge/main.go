package main

import (
    "bufio"
    "encoding/json"
    "flag"
    "fmt"
    "os"
    "os/exec"
    "os/signal"
    "strings"
    "syscall"
    "time"
)

var version = "4.0.0"

// ScanResult holds all scan results
type ScanResult struct {
    Target         string
    Subdomains     []string
    LiveHosts      []string
    URLs           []string
    Vulnerabilities []string
    StartTime      time.Time
    EndTime        time.Time
}

func main() {
    var target string
    var deep bool
    var threads int
    var outputDir string
    var showVersion bool

    flag.StringVar(&target, "t", "", "Target domain (required)")
    flag.StringVar(&target, "target", "", "Target domain (required)")
    flag.BoolVar(&deep, "deep", false, "Deep scan mode (enables port scanning, screenshots)")
    flag.IntVar(&threads, "T", 50, "Number of threads for concurrent scanning")
    flag.StringVar(&outputDir, "o", "", "Custom output directory")
    flag.BoolVar(&showVersion, "version", false, "Show version")
    flag.Parse()

    if showVersion {
        fmt.Printf("ReconForge v%s\n", version)
        os.Exit(0)
    }

    if target == "" {
        fmt.Println("\n❌ Error: Target required")
        fmt.Println("\nUsage: ./reconforge -t example.com [OPTIONS]")
        fmt.Println("\nOptions:")
        fmt.Println("  -t, --target     Target domain (required)")
        fmt.Println("  -d, --deep       Deep scan mode")
        fmt.Println("  -T, --threads    Number of threads (default: 50)")
        fmt.Println("  -o, --output     Custom output directory")
        fmt.Println("  --version        Show version")
        fmt.Println("\nExamples:")
        fmt.Println("  ./reconforge -t example.com")
        fmt.Println("  ./reconforge -t example.com --deep")
        fmt.Println("  ./reconforge -t example.com -T 100")
        os.Exit(1)
    }

    // Clean target
    target = strings.TrimPrefix(target, "https://")
    target = strings.TrimPrefix(target, "http://")
    target = strings.TrimPrefix(target, "www.")
    target = strings.Split(target, "/")[0]

    // Create output directory
    if outputDir == "" {
        timestamp := time.Now().Format("20060102_150405")
        outputDir = fmt.Sprintf("recon_%s_%s", target, timestamp)
    }
    os.MkdirAll(outputDir, 0755)

    // Setup signal handling for graceful shutdown
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    go func() {
        <-sigChan
        fmt.Println("\n\n⚠️  Interrupt received, cleaning up...")
        os.Exit(0)
    }()

    // Print banner
    printBanner()
    fmt.Printf("\n  🎯 Target: %s\n", target)
    fmt.Printf("  🔧 Mode: %s\n", map[bool]string{true: "Deep Scan", false: "Standard Scan"}[deep])
    fmt.Printf("  ⚡ Threads: %d\n", threads)
    fmt.Printf("  📁 Output: %s\n", outputDir)
    fmt.Println("\n" + strings.Repeat("─", 60))

    result := &ScanResult{
        Target:    target,
        StartTime: time.Now(),
    }

    // =============================================
    // PHASE 1: Subdomain Enumeration
    // =============================================
    fmt.Println("\n📡 PHASE 1: Subdomain Enumeration")
    result.Subdomains = enumerateSubdomains(target)
    fmt.Printf("  ✅ Found %d unique subdomains\n", len(result.Subdomains))
    saveToFile(outputDir+"/subdomains.txt", result.Subdomains)

    // =============================================
    // PHASE 2: Live Host Detection
    // =============================================
    if len(result.Subdomains) > 0 {
        fmt.Println("\n🌐 PHASE 2: Live Host Detection")
        result.LiveHosts = checkLiveHosts(result.Subdomains, threads)
        fmt.Printf("  ✅ Found %d live hosts\n", len(result.LiveHosts))
        saveToFile(outputDir+"/live_hosts.txt", result.LiveHosts)
    }

    // =============================================
    // PHASE 3: URL Discovery
    // =============================================
    fmt.Println("\n🔗 PHASE 3: URL Discovery")
    result.URLs = discoverURLs(target)
    fmt.Printf("  ✅ Found %d unique URLs\n", len(result.URLs))
    saveToFile(outputDir+"/urls.txt", result.URLs)

    // =============================================
    // PHASE 4: Vulnerability Scan (if deep mode)
    // =============================================
    if deep && len(result.LiveHosts) > 0 {
        fmt.Println("\n⚠️  PHASE 4: Vulnerability Assessment")
        result.Vulnerabilities = scanVulnerabilities(result.LiveHosts)
        fmt.Printf("  ✅ Found %d potential vulnerabilities\n", len(result.Vulnerabilities))
        saveToFile(outputDir+"/vulnerabilities.txt", result.Vulnerabilities)
    }

    result.EndTime = time.Now()

    // =============================================
    // Generate Reports
    // =============================================
    fmt.Println("\n📊 Generating Reports...")
    generateJSONReport(outputDir, result)
    generateHTMLReport(outputDir, result)

    // Final summary
    fmt.Println("\n" + strings.Repeat("─", 60))
    fmt.Println("\n✅ SCAN COMPLETED SUCCESSFULLY!")
    fmt.Printf("  📁 Results saved in: %s\n", outputDir)
    fmt.Printf("  ⏱️  Duration: %v\n", result.EndTime.Sub(result.StartTime))
    fmt.Printf("  🌐 Subdomains: %d\n", len(result.Subdomains))
    fmt.Printf("  🖥️  Live Hosts: %d\n", len(result.LiveHosts))
    fmt.Printf("  🔗 URLs: %d\n", len(result.URLs))
    if deep {
        fmt.Printf("  ⚠️  Vulnerabilities: %d (verify manually)\n", len(result.Vulnerabilities))
    }
    fmt.Println()
}

func printBanner() {
    fmt.Print("\033[36m")
    fmt.Println("╔══════════════════════════════════════════════════════════════════╗")
    fmt.Println("║                                                                  ║")
    fmt.Println("║   ██████╗ ███████╗ ██████╗ ██████╗ ███╗   ██╗███████╗ ██████╗    ║")
    fmt.Println("║   ██╔══██╗██╔════╝██╔════╝██╔═══██╗████╗  ██║██╔════╝██╔═══██╗   ║")
    fmt.Println("║   ██████╔╝█████╗  ██║     ██║   ██║██╔██╗ ██║█████╗  ██║   ██║   ║")
    fmt.Println("║   ██╔══██╗██╔══╝  ██║     ██║   ██║██║╚██╗██║██╔══╝  ██║   ██║   ║")
    fmt.Println("║   ██║  ██║███████╗╚██████╗╚██████╔╝██║ ╚████║██║     ╚██████╔╝   ║")
    fmt.Println("║   ╚═╝  ╚═╝╚══════╝ ╚═════╝ ╚═════╝ ╚═╝  ╚═══╝╚═╝      ╚═════╝    ║")
    fmt.Println("║                                                                  ║")
    fmt.Println("║                    ENTERPRISE RECONNAISSANCE                    ║")
    fmt.Println("║                         CODED BY : UMAR RUMAN                    ║")
    fmt.Println("║                      [ CYBER EX STUDY ]                         ║")
    fmt.Println("║                                                                  ║")
    fmt.Println("║              ⚠️  FOR AUTHORIZED USE ONLY  ⚠️                      ║")
    fmt.Println("║                                                                  ║")
    fmt.Println("║         📱 INSTAGRAM : @CYBER_EX_3697                            ║")
    fmt.Println("║         ▶️ YOUTUBE  : /@CyberEX3697                              ║")
    fmt.Println("║         💻 GITHUB  : /cyber-ex-3697                             ║")
    fmt.Println("║                                                                  ║")
    fmt.Println("╚══════════════════════════════════════════════════════════════════╝")
    fmt.Print("\033[0m")
}

func enumerateSubdomains(target string) []string {
    subdomains := make(map[string]bool)

    // Tool 1: subfinder
    fmt.Print("  🔍 Running subfinder... ")
    cmd := exec.Command("subfinder", "-d", target, "-silent")
    output, err := cmd.Output()
    if err == nil {
        for _, s := range strings.Split(string(output), "\n") {
            if s != "" && strings.Contains(s, target) {
                subdomains[s] = true
            }
        }
        fmt.Println("✅")
    } else {
        fmt.Println("❌ (not installed)")
    }

    // Tool 2: assetfinder
    fmt.Print("  🔍 Running assetfinder... ")
    cmd = exec.Command("assetfinder", "--subs-only", target)
    output, err = cmd.Output()
    if err == nil {
        for _, s := range strings.Split(string(output), "\n") {
            if s != "" && strings.Contains(s, target) {
                subdomains[s] = true
            }
        }
        fmt.Println("✅")
    } else {
        fmt.Println("❌ (not installed)")
    }

    // Tool 3: findomain
    fmt.Print("  🔍 Running findomain... ")
    cmd = exec.Command("findomain", "-t", target, "-q")
    output, err = cmd.Output()
    if err == nil {
        for _, s := range strings.Split(string(output), "\n") {
            if s != "" && strings.Contains(s, target) {
                subdomains[s] = true
            }
        }
        fmt.Println("✅")
    } else {
        fmt.Println("❌ (not installed)")
    }

    // Convert map to slice
    result := make([]string, 0, len(subdomains))
    for s := range subdomains {
        result = append(result, s)
    }
    return result
}

func checkLiveHosts(subdomains []string, threads int) []string {
    if len(subdomains) == 0 {
        return []string{}
    }

    // Create temp file
    tempFile := "/tmp/subdomains.txt"
    content := strings.Join(subdomains, "\n")
    os.WriteFile(tempFile, []byte(content), 0644)

    fmt.Print("  🌐 Running httpx... ")
    cmd := exec.Command("httpx", "-l", tempFile, "-silent", "-status-code", "-threads", fmt.Sprintf("%d", threads))
    output, err := cmd.Output()
    if err != nil {
        fmt.Println("❌")
        return []string{}
    }
    fmt.Println("✅")

    var live []string
    scanner := bufio.NewScanner(strings.NewReader(string(output)))
    for scanner.Scan() {
        line := scanner.Text()
        if strings.Contains(line, "[200]") || strings.Contains(line, "[301]") || strings.Contains(line, "[302]") {
            parts := strings.Fields(line)
            if len(parts) > 0 {
                live = append(live, parts[0])
            }
        }
    }
    return live
}

func discoverURLs(target string) []string {
    fmt.Print("  🔗 Running gau... ")
    cmd := exec.Command("gau", "--subs", target)
    output, err := cmd.Output()
    if err != nil {
        fmt.Println("❌")
        return []string{}
    }
    fmt.Println("✅")

    urls := make(map[string]bool)
    for _, u := range strings.Split(string(output), "\n") {
        if u != "" {
            urls[u] = true
        }
    }

    result := make([]string, 0, len(urls))
    for u := range urls {
        result = append(result, u)
    }
    return result
}

func scanVulnerabilities(hosts []string) []string {
    if len(hosts) == 0 {
        return []string{}
    }

    tempFile := "/tmp/hosts.txt"
    content := strings.Join(hosts, "\n")
    os.WriteFile(tempFile, []byte(content), 0644)

    fmt.Print("  ⚠️  Running nuclei... ")
    cmd := exec.Command("nuclei", "-l", tempFile, "-silent", "-severity", "critical,high", "-rate-limit", "10")
    output, err := cmd.Output()
    if err != nil {
        fmt.Println("❌")
        return []string{}
    }
    fmt.Println("✅")

    var findings []string
    scanner := bufio.NewScanner(strings.NewReader(string(output)))
    for scanner.Scan() {
        line := scanner.Text()
        if line != "" {
            findings = append(findings, line)
        }
    }
    return findings
}

func saveToFile(filename string, data []string) {
    content := strings.Join(data, "\n")
    os.WriteFile(filename, []byte(content), 0644)
}

func generateJSONReport(outputDir string, result *ScanResult) {
    report := map[string]interface{}{
        "target":         result.Target,
        "start_time":     result.StartTime,
        "end_time":       result.EndTime,
        "duration":       result.EndTime.Sub(result.StartTime).String(),
        "subdomains":     len(result.Subdomains),
        "live_hosts":     len(result.LiveHosts),
        "urls":           len(result.URLs),
        "vulnerabilities": len(result.Vulnerabilities),
    }
    data, _ := json.MarshalIndent(report, "", "  ")
    os.WriteFile(outputDir+"/report.json", data, 0644)
    fmt.Printf("  ✅ JSON report: %s/report.json\n", outputDir)
}

func generateHTMLReport(outputDir string, result *ScanResult) {
    htmlContent := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <title>ReconForge Report - %s</title>
    <meta charset="UTF-8">
    <style>
        body { background: #0a0e27; color: #e0e0e0; font-family: monospace; margin: 40px; }
        .container { max-width: 1200px; margin: auto; }
        .header { text-align: center; border-bottom: 2px solid #00ff88; padding-bottom: 20px; }
        .summary { display: grid; grid-template-columns: repeat(4, 1fr); gap: 20px; margin: 30px 0; }
        .card { background: #1a1f3a; padding: 20px; border-radius: 10px; border-left: 4px solid #00ff88; }
        .number { font-size: 2em; font-weight: bold; color: #00ff88; }
        .section { background: #0f1433; padding: 20px; border-radius: 10px; margin: 20px 0; }
        h1, h2 { color: #00ff88; }
        pre { background: #1a1f3a; padding: 15px; overflow-x: auto; border-radius: 5px; }
        .footer { text-align: center; padding: 20px; color: #666; border-top: 1px solid #2a2f4a; margin-top: 30px; }
    </style>
</head>
<body>
<div class="container">
    <div class="header">
        <h1>🔍 ReconForge Report</h1>
        <p>%s | %s</p>
    </div>
    <div class="summary">
        <div class="card"><h3>🌐 Subdomains</h3><div class="number">%d</div></div>
        <div class="card"><h3>✅ Live Hosts</h3><div class="number">%d</div></div>
        <div class="card"><h3>🔗 URLs</h3><div class="number">%d</div></div>
        <div class="card"><h3>⚠️ Findings</h3><div class="number">%d</div></div>
    </div>
    <div class="section">
        <h2>📋 Live Hosts</h2>
        <pre>%s</pre>
    </div>
    <div class="footer">
        <p>Generated by ReconForge v4.0 | Authorized Use Only</p>
    </div>
</div>
</body>
</html>`,
        result.Target,
        result.Target,
        time.Now().Format("2006-01-02 15:04:05"),
        len(result.Subdomains),
        len(result.LiveHosts),
        len(result.URLs),
        len(result.Vulnerabilities),
        strings.Join(result.LiveHosts, "\n"),
    )
    os.WriteFile(outputDir+"/report.html", []byte(htmlContent), 0644)
    fmt.Printf("  ✅ HTML report: %s/report.html\n", outputDir)
}
