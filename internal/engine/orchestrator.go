package engine

import (
    "bufio"
    "context"
    "fmt"
    "os"
    "os/exec"
    "strings"
    "time"
    
    "reconforge/internal/config"
    "reconforge/internal/logger"
)

type Orchestrator struct {
    config *config.Config
    logger *logger.Logger
    state  *ScanState
}

func NewOrchestrator(cfg *config.Config, log *logger.Logger, state *ScanState) *Orchestrator {
    return &Orchestrator{
        config: cfg,
        logger: log,
        state:  state,
    }
}

func (o *Orchestrator) RunPhase1(ctx context.Context, target string, result *ScanResult) error {
    o.logger.Info("Phase 1: Subdomain Enumeration")
    
    subdomains := make(map[string]bool)
    
    // subfinder
    o.logger.Debug("Running subfinder...")
    cmd := exec.CommandContext(ctx, "subfinder", "-d", target, "-silent")
    output, err := cmd.Output()
    if err == nil {
        for _, s := range strings.Split(string(output), "\n") {
            if s != "" && strings.Contains(s, target) {
                subdomains[s] = true
            }
        }
    }
    
    // assetfinder
    o.logger.Debug("Running assetfinder...")
    cmd = exec.CommandContext(ctx, "assetfinder", "--subs-only", target)
    output, err = cmd.Output()
    if err == nil {
        for _, s := range strings.Split(string(output), "\n") {
            if s != "" && strings.Contains(s, target) {
                subdomains[s] = true
            }
        }
    }
    
    // Convert map to slice
    for s := range subdomains {
        result.Subdomains = append(result.Subdomains, s)
    }
    
    o.logger.Success(fmt.Sprintf("Found %d subdomains", len(result.Subdomains)))
    return nil
}

func (o *Orchestrator) RunPhase2(ctx context.Context, subdomains []string, result *ScanResult) error {
    o.logger.Info("Phase 2: Live Host Detection")
    
    if len(subdomains) == 0 {
        return nil
    }
    
    // Create temp file
    tempFile := "/tmp/subdomains.txt"
    content := strings.Join(subdomains, "\n")
    os.WriteFile(tempFile, []byte(content), 0644)
    
    cmd := exec.CommandContext(ctx, "httpx", "-l", tempFile, "-silent", "-status-code", "-threads", fmt.Sprintf("%d", o.config.Scan.Threads))
    output, err := cmd.Output()
    if err != nil {
        return err
    }
    
    scanner := bufio.NewScanner(strings.NewReader(string(output)))
    for scanner.Scan() {
        line := scanner.Text()
        if strings.Contains(line, "[200]") || strings.Contains(line, "[301]") || strings.Contains(line, "[302]") {
            parts := strings.Fields(line)
            if len(parts) > 0 {
                result.LiveHosts = append(result.LiveHosts, parts[0])
            }
        }
    }
    
    o.logger.Success(fmt.Sprintf("Found %d live hosts", len(result.LiveHosts)))
    return nil
}

func (o *Orchestrator) RunPhase3(ctx context.Context, target string, result *ScanResult) error {
    o.logger.Info("Phase 3: URL Discovery")
    
    cmd := exec.CommandContext(ctx, "gau", "--subs", target)
    output, err := cmd.Output()
    if err != nil {
        return err
    }
    
    urls := make(map[string]bool)
    for _, u := range strings.Split(string(output), "\n") {
        if u != "" {
            urls[u] = true
        }
    }
    
    for u := range urls {
        result.URLs = append(result.URLs, u)
    }
    
    o.logger.Success(fmt.Sprintf("Found %d URLs", len(result.URLs)))
    return nil
}

func (o *Orchestrator) RunPhase4(ctx context.Context, liveHosts []string, result *ScanResult) error {
    o.logger.Info("Phase 4: Vulnerability Assessment")
    
    if len(liveHosts) == 0 {
        return nil
    }
    
    tempFile := "/tmp/hosts.txt"
    content := strings.Join(liveHosts, "\n")
    os.WriteFile(tempFile, []byte(content), 0644)
    
    cmd := exec.CommandContext(ctx, "nuclei", "-l", tempFile, "-silent", "-severity", "critical,high", "-rate-limit", "10")
    output, err := cmd.Output()
    if err != nil {
        return err
    }
    
    scanner := bufio.NewScanner(strings.NewReader(string(output)))
    for scanner.Scan() {
        line := scanner.Text()
        if line != "" {
            result.Vulnerabilities = append(result.Vulnerabilities, line)
        }
    }
    
    o.logger.Success(fmt.Sprintf("Found %d potential vulnerabilities", len(result.Vulnerabilities)))
    return nil
}
