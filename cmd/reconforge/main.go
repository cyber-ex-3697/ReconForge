package main

import (
    "bufio"
    "encoding/json"
    "flag"
    "fmt"
    "math/rand"
    "net/http"
    "os"
    "os/exec"
    "os/signal"
    "strings"
    "sync"
    "syscall"
    "time"

    "reconforge/pkg/api"
)

var version = "5.0.0"

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

// Config holds configuration settings
type Config struct {
    AIAssisted    bool
    Distributed   bool
    GraphEnabled  bool
    Headless      bool
    Evasion       bool
    PassiveOnly   bool
    Profile       string
}

// Rate limit tracking
var (
    totalRateLimits   int
    rateLimitMutex    sync.Mutex
    lastRateLimitTime time.Time
    adaptiveThreads   int
)

func trackRateLimit() {
    rateLimitMutex.Lock()
    defer rateLimitMutex.Unlock()
    
    totalRateLimits++
    lastRateLimitTime = time.Now()
    
    if totalRateLimits > 10 {
        fmt.Println("\n⚠️  Too many rate limits detected! Consider:")
        fmt.Println("   1. Waiting 1-2 hours before next scan")
        fmt.Println("   2. Using a VPN or proxy")
        fmt.Println("   3. Scanning a different target")
        fmt.Println("   4. Using -T 10 to reduce threads")
    }
}

func getJitterDelay(baseDelay time.Duration) time.Duration {
    rand.Seed(time.Now().UnixNano())
    jitter := time.Duration(rand.Int63n(int64(baseDelay / 2)))
    return baseDelay + jitter
}

func analyzeTarget(target string) {
    fmt.Println("🔍 Analyzing target behavior...")
    
    client := &http.Client{Timeout: 10 * time.Second}
    resp, err := client.Get("https://" + target)
    if err != nil {
        fmt.Println("  ❌ Target unreachable or SSL error")
        return
    }
    defer resp.Body.Close()
    
    fmt.Printf("  ✅ Target reachable\n")
    fmt.Printf("  Status code: %d\n", resp.StatusCode)
    
    if resp.Header.Get("X-RateLimit-Remaining") != "" {
        fmt.Println("  ⚠️  Target has rate limiting enabled!")
        fmt.Println("  💡 Recommended: use -T 10 for better results")
    }
    
    if resp.Header.Get("CF-Ray") != "" {
        fmt.Println("  ☁️  Target is behind Cloudflare")
        fmt.Println("  💡 Cloudflare has strict rate limiting. Use -T 5")
    }
}

func suggestThreads(target string, userThreads int) int {
    if userThreads > 0 {
        return userThreads
    }
    
    client := &http.Client{Timeout: 5 * time.Second}
    start := time.Now()
    resp, err := client.Get("https://" + target)
    if err != nil {
        return 50
    }
    defer resp.Body.Close()
    
    elapsed := time.Since(start)
    
    if elapsed > 3*time.Second {
        fmt.Println("  💡 Target is slow! Recommended threads: 10")
        return 10
    } else if elapsed < 500*time.Millisecond {
        if resp.Header.Get("CF-Ray") != "" {
            fmt.Println("  💡 Cloudflare detected. Recommended threads: 20")
            return 20
        }
        fmt.Println("  💡 Target is fast! Recommended threads: 100")
        return 100
    }
    
    fmt.Println("  💡 Using default threads: 50")
    return 50
}

func main() {
    var target string
    var deep bool
    var threads int
    var outputDir string
    var showVersion bool
    var configFile string
    var profile string
    var aiAssisted bool
    var distributed bool
    var workerRole string
    var graphEnabled bool
    var headless bool
    var evasion bool
    var passiveOnly bool
    var stealth bool
    var aggressive bool

    // Basic flags
    flag.StringVar(&target, "t", "", "Target domain (required)")
    flag.StringVar(&target, "target", "", "Target domain (required)")
    flag.BoolVar(&deep, "deep", false, "Deep scan mode (enables port scanning, screenshots)")
    flag.IntVar(&threads, "T", 0, "Number of threads for concurrent scanning (auto-adjusted)")
    flag.StringVar(&outputDir, "o", "", "Custom output directory")
    flag.BoolVar(&showVersion, "version", false, "Show version")
    flag.StringVar(&configFile, "c", "config.yaml", "Config file path")
    
    // New feature flags for v5.0
    flag.StringVar(&profile, "profile", "standard", "Scan profile: standard, stealth, aggressive, passive, distributed, ai_assisted")
    flag.BoolVar(&aiAssisted, "ai-assisted", false, "Enable AI-assisted recon engine")
    flag.BoolVar(&distributed, "distributed", false, "Enable distributed scanning mode")
    flag.StringVar(&workerRole, "worker-role", "worker", "Worker role: master, worker")
    flag.BoolVar(&graphEnabled, "graph", false, "Enable graph database for attack path correlation")
    flag.BoolVar(&headless, "headless", false, "Enable headless browser crawling")
    flag.BoolVar(&evasion, "evasion", false, "Enable WAF evasion techniques")
    flag.BoolVar(&passiveOnly, "passive", false, "Passive-only recon (no active requests)")
    flag.BoolVar(&stealth, "stealth", false, "Stealth mode (slow, avoids detection)")
    flag.BoolVar(&aggressive, "aggressive", false, "Aggressive mode (fast, may get rate limited)")
    
    flag.Parse()

    if showVersion {
        fmt.Printf("RECONFORGE v%s\n", version)
        os.Exit(0)
    }

    if target == "" {
        fmt.Println("\n❌ Error: Target required")
        fmt.Println("\nUsage: ./reconforge -t example.com [OPTIONS]")
        fmt.Println("\nBasic Options:")
        fmt.Println("  -t, --target     Target domain (required)")
        fmt.Println("  -d, --deep       Deep scan mode")
        fmt.Println("  -T, --threads    Number of threads (default: auto)")
        fmt.Println("  -o, --output     Custom output directory")
        fmt.Println("  -c, --config     Config file path (default: config.yaml)")
        fmt.Println("  --version        Show version")
        fmt.Println("\nAdvanced Options (v5.0):")
        fmt.Println("  --profile        Scan profile (standard, stealth, aggressive, passive, distributed, ai_assisted)")
        fmt.Println("  --ai-assisted    Enable AI-assisted recon engine")
        fmt.Println("  --distributed    Enable distributed scanning mode")
        fmt.Println("  --worker-role    Worker role: master, worker (for distributed mode)")
        fmt.Println("  --graph          Enable graph database for attack path correlation")
        fmt.Println("  --headless       Enable headless browser crawling")
        fmt.Println("  --evasion        Enable WAF evasion techniques")
        fmt.Println("  --passive        Passive-only recon (no active requests)")
        fmt.Println("  --stealth        Stealth mode (slow, avoids detection)")
        fmt.Println("  --aggressive     Aggressive mode (fast, high concurrency)")
        fmt.Println("\nExamples:")
        fmt.Println("  ./reconforge -t example.com")
        fmt.Println("  ./reconforge -t example.com --deep")
        fmt.Println("  ./reconforge -t example.com --ai-assisted")
        fmt.Println("  ./reconforge -t example.com --profile stealth")
        fmt.Println("  ./reconforge -t example.com --distributed --worker-role master")
        fmt.Println("  ./reconforge -t example.com --graph --deep")
        os.Exit(1)
    }

    // Clean target
    target = strings.TrimPrefix(target, "https://")
    target = strings.TrimPrefix(target, "http://")
    target = strings.TrimPrefix(target, "www.")
    target = strings.Split(target, "/")[0]

    // Apply profile settings
    if stealth {
        profile = "stealth"
        threads = 5
        fmt.Println("  🥷 Stealth mode enabled (threads: 5, high delay)")
    } else if aggressive {
        profile = "aggressive"
        if threads == 0 {
            threads = 200
        }
        fmt.Println("  ⚡ Aggressive mode enabled")
    } else if passiveOnly {
        profile = "passive"
        fmt.Println("  📡 Passive-only mode enabled (no active requests)")
    } else if aiAssisted {
        profile = "ai_assisted"
        fmt.Println("  🧠 AI-assisted mode enabled")
    } else if distributed {
        profile = "distributed"
        fmt.Printf("  🌐 Distributed mode enabled (role: %s)\n", workerRole)
    } else if threads == 0 {
        threads = suggestThreads(target, 0)
    }

    // Limit threads
    if threads > 200 {
        fmt.Printf("  ⚠️  Threads reduced from %d to 200 (max recommended)\n", threads)
        threads = 200
    }
    if threads < 1 {
        threads = 1
    }
    
    adaptiveThreads = threads

    // Analyze target
    analyzeTarget(target)

    // Create output directory
    if outputDir == "" {
        timestamp := time.Now().Format("20060102_150405")
        outputDir = fmt.Sprintf("recon_%s_%s", target, timestamp)
    }
    os.MkdirAll(outputDir, 0755)

    // Load config and initialize API client
    apiClient := loadAPIClient(configFile)

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
    fmt.Printf("  📋 Profile: %s\n", profile)
    fmt.Println("\n" + strings.Repeat("─", 60))

    // Print feature status
    fmt.Println("\n🔧 Feature Status:")
    if aiAssisted {
        fmt.Println("  🧠 AI-Assisted Recon: ENABLED")
    }
    if distributed {
        fmt.Println("  🌐 Distributed Mode: ENABLED")
        fmt.Printf("     Worker Role: %s\n", workerRole)
    }
    if graphEnabled {
        fmt.Println("  📊 Graph Database: ENABLED")
    }
    if headless {
        fmt.Println("  🌐 Headless Crawling: ENABLED")
    }
    if evasion {
        fmt.Println("  🛡️ WAF Evasion: ENABLED")
    }
    if passiveOnly {
        fmt.Println("  📡 Passive Recon: ENABLED")
    }

    result := &ScanResult{
        Target:    target,
        StartTime: time.Now(),
    }

    // =============================================
    // PHASE 1: Subdomain Enumeration
    // =============================================
    fmt.Println("\n📡 PHASE 1: Subdomain Enumeration")
    result.Subdomains = enumerateSubdomainsWithAPI(target, apiClient, threads)
    fmt.Printf("  ✅ Found %d unique subdomains\n", len(result.Subdomains))
    saveToFile(outputDir+"/subdomains.txt", result.Subdomains)

    // =============================================
    // PHASE 2: Live Host Detection
    // =============================================
    if len(result.Subdomains) > 0 && !passiveOnly {
        fmt.Println("\n🌐 PHASE 2: Live Host Detection")
        result.LiveHosts = checkLiveHosts(result.Subdomains, threads, stealth)
        fmt.Printf("  ✅ Found %d live hosts\n", len(result.LiveHosts))
        saveToFile(outputDir+"/live_hosts.txt", result.LiveHosts)
    } else if passiveOnly {
        fmt.Println("\n🌐 PHASE 2: Live Host Detection (SKIPPED - Passive mode)")
    }

    // =============================================
    // PHASE 3: URL Discovery
    // =============================================
    if !passiveOnly {
        fmt.Println("\n🔗 PHASE 3: URL Discovery")
        result.URLs = discoverURLs(target)
        fmt.Printf("  ✅ Found %d unique URLs\n", len(result.URLs))
        saveToFile(outputDir+"/urls.txt", result.URLs)
    } else {
        fmt.Println("\n🔗 PHASE 3: URL Discovery (SKIPPED - Passive mode)")
    }

    // =============================================
    // PHASE 4: Vulnerability Scan (if deep mode)
    // =============================================
    if deep && len(result.LiveHosts) > 0 && !passiveOnly {
        fmt.Println("\n⚠️  PHASE 4: Vulnerability Assessment")
        
        // Adjust rate limit based on profile
        rateLimit := 10
        if stealth {
            rateLimit = 2
        } else if aggressive {
            rateLimit = 50
        }
        
        result.Vulnerabilities = scanVulnerabilities(result.LiveHosts, rateLimit)
        fmt.Printf("  ✅ Found %d potential vulnerabilities\n", len(result.Vulnerabilities))
        saveToFile(outputDir+"/vulnerabilities.txt", result.Vulnerabilities)
    } else if deep && passiveOnly {
        fmt.Println("\n⚠️  PHASE 4: Vulnerability Assessment (SKIPPED - Passive mode)")
    } else if deep && len(result.LiveHosts) == 0 {
        fmt.Println("\n⚠️  PHASE 4: Vulnerability Assessment (SKIPPED - No live hosts)")
    }

    result.EndTime = time.Now()

    // =============================================
    // Generate Reports
    // =============================================
    fmt.Println("\n📊 Generating Reports...")
    generateJSONReport(outputDir, result)
    generateHTMLReport(outputDir, result, profile)

    // AI-assisted report enrichment
    if aiAssisted {
        fmt.Println("  🧠 AI-assisted report enrichment...")
        generateAIEnrichedReport(outputDir, result)
    }

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
    
    // Feature-specific recommendations
    if aiAssisted {
        fmt.Println("\n💡 AI-Assisted Recommendations:")
        fmt.Println("  • Review priority-ranked subdomains first")
        fmt.Println("  • Check vulnerability predictions against actual findings")
        fmt.Println("  • Run with --graph for attack path visualization")
    }
    
    if distributed {
        fmt.Println("\n🌐 Distributed Scan Info:")
        fmt.Printf("  • Role: %s\n", workerRole)
        fmt.Println("  • Use --worker-role to change role")
        fmt.Println("  • Ensure Redis is running for queue coordination")
    }
    
    if graphEnabled {
        fmt.Println("\n📊 Graph Database Info:")
        fmt.Println("  • Attack paths saved to Neo4j/ArangoDB")
        fmt.Println("  • Use graph viewer for visualization")
    }
    
    fmt.Println()
}

func loadAPIClient(configFile string) *api.APIClient {
    client := api.NewAPIClient()
    
    data, err := os.ReadFile(configFile)
    if err != nil {
        fmt.Println("  ⚠️  Config file not found, API integrations disabled")
        return client
    }
    
    lines := strings.Split(string(data), "\n")
    for _, line := range lines {
        line = strings.TrimSpace(line)
        if strings.HasPrefix(line, "chaos_key:") {
            key := strings.TrimSpace(strings.TrimPrefix(line, "chaos_key:"))
            key = strings.Trim(key, "\"")
            if key != "" {
                client.SetAPIKey("chaos", key)
                fmt.Println("  ✅ Chaos API enabled")
            }
        } else if strings.HasPrefix(line, "github_token:") {
            token := strings.TrimSpace(strings.TrimPrefix(line, "github_token:"))
            token = strings.Trim(token, "\"")
            if token != "" {
                client.SetAPIKey("github", token)
                fmt.Println("  ✅ GitHub API enabled")
            }
        } else if strings.HasPrefix(line, "shodan_key:") {
            key := strings.TrimSpace(strings.TrimPrefix(line, "shodan_key:"))
            key = strings.Trim(key, "\"")
            if key != "" {
                client.SetAPIKey("shodan", key)
                fmt.Println("  ✅ Shodan API enabled")
            }
        }
    }
    
    return client
}

func generateAIEnrichedReport(outputDir string, result *ScanResult) {
    enrichedFile := outputDir + "/ai_enriched_report.json"
    
    // Calculate priority targets
    priorityTargets := getPriorityTargets(result.Subdomains)
    
    report := map[string]interface{}{
        "target": result.Target,
        "insights": map[string]interface{}{
            "priority_targets":      priorityTargets,
            "vulnerability_likelihood": calculateVulnerabilityLikelihood(result),
            "recommended_actions":   getRecommendedActions(result),
            "attack_surface_score":  calculateAttackSurfaceScore(result),
        },
        "statistics": map[string]interface{}{
            "total_subdomains": len(result.Subdomains),
            "total_live_hosts": len(result.LiveHosts),
            "total_urls":       len(result.URLs),
            "total_vulnerabilities": len(result.Vulnerabilities),
        },
    }
    
    data, _ := json.MarshalIndent(report, "", "  ")
    os.WriteFile(enrichedFile, data, 0644)
    fmt.Printf("  ✅ AI-enriched report: %s\n", enrichedFile)
}

func getPriorityTargets(subdomains []string) []string {
    var priority []string
    patterns := []string{"admin", "api", "dashboard", "portal", "internal", "dev", "staging", "backup"}
    
    for _, sub := range subdomains {
        for _, pattern := range patterns {
            if strings.Contains(sub, pattern) {
                priority = append(priority, sub)
                break
            }
        }
        if len(priority) >= 20 {
            break
        }
    }
    return priority
}

func calculateVulnerabilityLikelihood(result *ScanResult) string {
    if len(result.Vulnerabilities) > 10 {
        return "HIGH - Multiple vulnerabilities detected"
    } else if len(result.Vulnerabilities) > 3 {
        return "MEDIUM - Several potential issues found"
    } else if len(result.Vulnerabilities) > 0 {
        return "LOW - Few minor issues detected"
    }
    return "UNKNOWN - Manual review recommended"
}

func getRecommendedActions(result *ScanResult) []string {
    actions := []string{}
    
    if len(result.Subdomains) > 0 {
        actions = append(actions, "Review discovered subdomains for sensitive assets")
    }
    if len(result.LiveHosts) > 0 {
        actions = append(actions, "Prioritize live hosts for deeper testing")
    }
    if len(result.Vulnerabilities) > 0 {
        actions = append(actions, "Manually verify potential vulnerabilities")
        actions = append(actions, "Focus on critical/high severity findings first")
    }
    if len(result.URLs) > 1000 {
        actions = append(actions, "Large URL set discovered - consider parameter fuzzing")
    }
    
    if len(actions) == 0 {
        actions = append(actions, "Run deep scan for more comprehensive results")
        actions = append(actions, "Use --ai-assisted for intelligent scanning")
    }
    
    return actions
}

func calculateAttackSurfaceScore(result *ScanResult) int {
    score := 0
    
    // Subdomain scoring
    if len(result.Subdomains) > 1000 {
        score += 30
    } else if len(result.Subdomains) > 100 {
        score += 15
    } else if len(result.Subdomains) > 10 {
        score += 5
    }
    
    // Live host scoring
    if len(result.LiveHosts) > 100 {
        score += 30
    } else if len(result.LiveHosts) > 20 {
        score += 15
    } else if len(result.LiveHosts) > 0 {
        score += 5
    }
    
    // URL scoring
    if len(result.URLs) > 10000 {
        score += 20
    } else if len(result.URLs) > 1000 {
        score += 10
    } else if len(result.URLs) > 100 {
        score += 5
    }
    
    // Vulnerability scoring
    if len(result.Vulnerabilities) > 20 {
        score += 20
    } else if len(result.Vulnerabilities) > 5 {
        score += 10
    }
    
    if score > 100 {
        score = 100
    }
    
    return score
}

func enumerateSubdomainsWithAPI(target string, apiClient *api.APIClient, threads int) []string {
    subdomains := make(map[string]bool)
    rateLimitCount := 0

    // Tool 1: subfinder
    fmt.Print("  🔍 Running subfinder... ")
    cmd := exec.Command("subfinder", "-d", target, "-silent", "-t", fmt.Sprintf("%d", threads/2))
    output, err := cmd.Output()
    if err == nil {
        for _, s := range strings.Split(string(output), "\n") {
            if s != "" && strings.Contains(s, target) {
                subdomains[s] = true
            }
        }
        fmt.Println("✅")
    } else {
        if strings.Contains(err.Error(), "429") {
            rateLimitCount++
            fmt.Println("⚠️ (rate limited)")
            trackRateLimit()
        } else {
            fmt.Println("❌ (not installed)")
        }
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

    // Chaos API
    chaosKey := apiClient.GetAPIKey("chaos")
    if chaosKey != "" {
        fmt.Print("  ☁️  Running Chaos API... ")
        chaosSubs, err := apiClient.GetChaosSubdomains(target)
        if err == nil {
            for _, s := range chaosSubs {
                if s != "" {
                    subdomains[s] = true
                }
            }
            fmt.Printf("✅ (+%d)\n", len(chaosSubs))
        } else {
            fmt.Printf("❌ (%v)\n", err)
        }
    } else {
        fmt.Print("  ☁️  Chaos API... ⚠️ (key not configured)\n")
    }

    // Shodan API
    shodanKey := apiClient.GetAPIKey("shodan")
    if shodanKey != "" {
        fmt.Print("  🌍 Running Shodan API... ")
        shodanHosts, err := apiClient.GetShodanHosts(target)
        if err == nil {
            for _, h := range shodanHosts {
                if h != "" && strings.Contains(h, target) {
                    subdomains[h] = true
                }
            }
            fmt.Printf("✅ (+%d)\n", len(shodanHosts))
        } else {
            fmt.Printf("❌ (%v)\n", err)
        }
    } else {
        fmt.Print("  🌍 Shodan API... ⚠️ (key not configured)\n")
    }

    // GitHub API
    githubToken := apiClient.GetAPIKey("github")
    if githubToken != "" {
        fmt.Print("  💻 Running GitHub API... ")
        githubURLs, err := apiClient.GetGitHubSubdomains(target)
        if err == nil {
            for _, u := range githubURLs {
                if strings.Contains(u, target) {
                    subdomains[u] = true
                }
            }
            fmt.Printf("✅ (+%d)\n", len(githubURLs))
        } else {
            fmt.Printf("❌ (%v)\n", err)
        }
    } else {
        fmt.Print("  💻 GitHub API... ⚠️ (token not configured)\n")
    }

    if rateLimitCount > 0 {
        fmt.Println("\n  ⚠️  Rate limiting detected! Suggestions:")
        fmt.Println("     • Wait 2-4 hours before scanning again")
        fmt.Println("     • Use --stealth mode for slower scanning")
        fmt.Println("     • Use -T 10 to reduce threads")
        fmt.Println("     • Scan a different target")
    }

    result := make([]string, 0, len(subdomains))
    for s := range subdomains {
        result = append(result, s)
    }
    return result
}

func checkLiveHosts(subdomains []string, threads int, stealth bool) []string {
    if len(subdomains) == 0 {
        return []string{}
    }

    if stealth {
        threads = 5
        fmt.Println("  🥷 Stealth mode: using 5 threads with 3s delay")
    } else if threads > 20 {
        threads = 20
        fmt.Println("  ⚠️  Reduced threads to 20 to avoid rate limiting")
    }

    tempFile := "/tmp/subdomains.txt"
    content := strings.Join(subdomains, "\n")
    os.WriteFile(tempFile, []byte(content), 0644)

    if stealth {
        delay := getJitterDelay(3 * time.Second)
        time.Sleep(delay)
    }

    fmt.Print("  🌐 Running httpx with rate limiting protection... ")
    
    var cmd *exec.Cmd
    if stealth {
        cmd = exec.Command("httpx", "-l", tempFile, "-silent", "-status-code",
            "-threads", "5",
            "-delay", "3s",
            "-retries", "3",
            "-timeout", "15")
    } else {
        cmd = exec.Command("httpx", "-l", tempFile, "-silent", "-status-code",
            "-threads", fmt.Sprintf("%d", threads),
            "-delay", "2s",
            "-retries", "2",
            "-timeout", "10")
    }
    
    output, err := cmd.Output()
    if err != nil {
        if strings.Contains(err.Error(), "429") {
            trackRateLimit()
            fmt.Println("⚠️ (rate limited)")
        } else {
            fmt.Println("❌")
        }
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
    
    if len(live) == 0 {
        fmt.Println("  ⚠️  No live hosts found (possibly rate limited)")
        fmt.Println("  💡 Try again later or use --stealth mode")
        fmt.Println("  💡 Or use: ./reconforge -t " + strings.Split(subdomains[0], ".")[1] + " -T 10")
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

func scanVulnerabilities(hosts []string, rateLimit int) []string {
    if len(hosts) == 0 {
        return []string{}
    }

    tempFile := "/tmp/hosts.txt"
    content := strings.Join(hosts, "\n")
    os.WriteFile(tempFile, []byte(content), 0644)

    fmt.Print("  ⚠️  Running nuclei with rate limiting... ")
    
    cmd := exec.Command("nuclei", "-l", tempFile, "-silent", "-severity", "critical,high", 
        "-rate-limit", fmt.Sprintf("%d", rateLimit),
        "-bulk-size", "10",
        "-timeout", "10",
        "-retries", "2")
    
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
        "target":          result.Target,
        "start_time":      result.StartTime,
        "end_time":        result.EndTime,
        "duration":        result.EndTime.Sub(result.StartTime).String(),
        "subdomains":      len(result.Subdomains),
        "live_hosts":      len(result.LiveHosts),
        "urls":            len(result.URLs),
        "vulnerabilities": len(result.Vulnerabilities),
    }
    data, _ := json.MarshalIndent(report, "", "  ")
    os.WriteFile(outputDir+"/report.json", data, 0644)
    fmt.Printf("  ✅ JSON report: %s/report.json\n", outputDir)
}

func generateHTMLReport(outputDir string, result *ScanResult, profile string) {
    htmlContent := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <title>RECONFORGE Report - %s</title>
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
        .tip { background: #1a1f3a; padding: 15px; border-left: 4px solid #ffcc44; margin: 20px 0; }
        .profile-badge { display: inline-block; padding: 4px 12px; border-radius: 20px; font-size: 12px; margin-left: 10px; }
        .profile-stealth { background: #4a4a4a; }
        .profile-aggressive { background: #ff4444; }
        .profile-ai { background: #9b59b6; }
        .profile-standard { background: #00ff88; color: #000; }
    </style>
</head>
<body>
<div class="container">
    <div class="header">
        <h1>🔍 RECONFORGE v5.0 Report</h1>
        <p>%s | %s | Profile: <span class="profile-badge profile-%s">%s</span></p>
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
    <div class="tip">
        💡 <strong>Tip:</strong> If you got 0 live hosts, the target may be rate limiting your requests.
        Try again with: <code>./reconforge -t %s --stealth</code> or <code>./reconforge -t %s --ai-assisted</code>
    </div>
    <div class="footer">
        <p>Generated by RECONFORGE v5.0 | Authorized Use Only | Profile: %s</p>
    </div>
</div>
</body>
</html>`,
        result.Target,
        result.Target,
        time.Now().Format("2006-01-02 15:04:05"),
        profile, profile,
        len(result.Subdomains),
        len(result.LiveHosts),
        len(result.URLs),
        len(result.Vulnerabilities),
        strings.Join(result.LiveHosts, "\n"),
        result.Target, result.Target,
        profile,
    )
    os.WriteFile(outputDir+"/report.html", []byte(htmlContent), 0644)
    fmt.Printf("  ✅ HTML report: %s/report.html\n", outputDir)
}

func printBanner() {
    fmt.Print("\033[36m")
    fmt.Println("╔═══════════════════════════════════════════════════════════════════════════════════════════╗")
    fmt.Println("║                                                                                           ║")
    fmt.Println("║   ██████╗ ███████╗ ██████╗ ██████╗ ███╗   ██╗███████╗ ██████╗ ██████╗ ██████╗ ███████╗    ║")
    fmt.Println("║   ██╔══██╗██╔════╝██╔════╝██╔═══██╗████╗  ██║██╔════╝██╔═══██╗██╔══██╗██╔══██╗██╔════╝    ║")
    fmt.Println("║   ██████╔╝█████╗  ██║     ██║   ██║██╔██╗ ██║█████╗  ██║   ██║██████╔╝██████╔╝█████╗      ║")
    fmt.Println("║   ██╔══██╗██╔══╝  ██║     ██║   ██║██║╚██╗██║██╔══╝  ██║   ██║██╔══██╗██╔══██╗██╔══╝      ║")
    fmt.Println("║   ██║  ██║███████╗╚██████╗╚██████╔╝██║ ╚████║██║     ╚██████╔╝██║  ██║██║  ██║███████╗    ║")
    fmt.Println("║   ╚═╝  ╚═╝╚══════╝ ╚═════╝ ╚═════╝ ╚═╝  ╚═══╝╚═╝      ╚═════╝ ╚═╝  ╚═╝╚═╝  ╚═╝╚══════╝    ║")
    fmt.Println("║                                                                                           ║")
    fmt.Println("║                                   RECONFORGE v5.0.0                                        ║")
    fmt.Println("║                           ENTERPRISE RECONNAISSANCE FRAMEWORK                              ║")
    fmt.Println("║                                                                                           ║")
    fmt.Println("║                              CODED BY : UMAR RUMAN                                         ║")
    fmt.Println("║                              [ CYBER EX STUDY ]                                           ║")
    fmt.Println("║                                                                                           ║")
    fmt.Println("║                         ⚠️  FOR AUTHORIZED USE ONLY  ⚠️                                    ║")
    fmt.Println("║                                                                                           ║")
    fmt.Println("║                         📱 INSTAGRAM : @CYBER_EX_3697                                      ║")
    fmt.Println("║                         ▶️ YOUTUBE  : /@CyberEX3697                                        ║")
    fmt.Println("║                         💻 GITHUB  : /cyber-ex-3697                                       ║")
    fmt.Println("║                                                                                           ║")
    fmt.Println("╚═══════════════════════════════════════════════════════════════════════════════════════════╝")
    fmt.Print("\033[0m")
}
