package main

import (
    "bufio"
    "encoding/json"
    "flag"
    "fmt"
    "os"
    "strings"
    "time"
)

// Session represents a captured session
type Session struct {
    ID        string                 `json:"id"`
    Name      string                 `json:"name"`
    Requests  []Request              `json:"requests"`
    CreatedAt time.Time              `json:"created_at"`
}

// Request represents an HTTP request
type Request struct {
    ID        string            `json:"id"`
    Method    string            `json:"method"`
    URL       string            `json:"url"`
    Headers   map[string]string `json:"headers"`
    Body      string            `json:"body"`
    Response  *Response         `json:"response"`
    Timestamp time.Time         `json:"timestamp"`
}

// Response represents an HTTP response
type Response struct {
    StatusCode int               `json:"status_code"`
    Headers    map[string]string `json:"headers"`
    Body       string            `json:"body"`
    Length     int               `json:"length"`
}

// AnalysisResult represents session analysis result
type AnalysisResult struct {
    TotalRequests   int            `json:"total_requests"`
    UniqueHosts     int            `json:"unique_hosts"`
    Methods         map[string]int `json:"methods"`
    StatusCodes     map[int]int    `json:"status_codes"`
    AverageBodySize float64        `json:"average_body_size"`
    TotalDuration   string         `json:"total_duration"`
}

func main() {
    var sessionPath string
    var format string
    var outputPath string
    
    flag.StringVar(&sessionPath, "session", "", "Session JSON file")
    flag.StringVar(&format, "format", "text", "Output format (text, json)")
    flag.StringVar(&outputPath, "output", "", "Output file (optional)")
    flag.Parse()
    
    fmt.Println("=== ReconForge Session Analyzer ===\n")
    
    if sessionPath == "" {
        fmt.Println("Usage: session_analyzer -session <session.json>")
        fmt.Println("\nOptions:")
        fmt.Println("  -session  Session JSON file to analyze")
        fmt.Println("  -format   Output format (text, json)")
        fmt.Println("  -output   Output file (optional)")
        os.Exit(1)
    }
    
    // Load session
    session, err := loadSession(sessionPath)
    if err != nil {
        fmt.Printf("[!] Error loading session: %v\n", err)
        os.Exit(1)
    }
    
    fmt.Printf("[✓] Loaded session: %s\n", session.Name)
    fmt.Printf("    Requests: %d\n", len(session.Requests))
    fmt.Printf("    Created: %s\n", session.CreatedAt.Format("2006-01-02 15:04:05"))
    
    // Analyze session
    result := analyzeSession(session)
    
    // Output results
    output := formatResult(result, format, session)
    
    if outputPath != "" {
        if err := os.WriteFile(outputPath, []byte(output), 0644); err != nil {
            fmt.Printf("[!] Error writing output: %v\n", err)
        } else {
            fmt.Printf("\n[✓] Results saved to: %s\n", outputPath)
        }
    } else {
        fmt.Print("\n" + output)
    }
}

func loadSession(path string) (*Session, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }
    
    var session Session
    if err := json.Unmarshal(data, &session); err != nil {
        return nil, err
    }
    
    return &session, nil
}

func analyzeSession(session *Session) *AnalysisResult {
    result := &AnalysisResult{
        TotalRequests: len(session.Requests),
        Methods:       make(map[string]int),
        StatusCodes:   make(map[int]int),
    }
    
    hosts := make(map[string]bool)
    var totalBodySize int
    
    for _, req := range session.Requests {
        // Count methods
        result.Methods[req.Method]++
        
        // Track hosts
        // Parse URL (simplified)
        parts := strings.Split(req.URL, "/")
        if len(parts) >= 3 {
            host := parts[2]
            hosts[host] = true
        }
        
        if req.Response != nil {
            result.StatusCodes[req.Response.StatusCode]++
            totalBodySize += req.Response.Length
        }
    }
    
    result.UniqueHosts = len(hosts)
    if result.TotalRequests > 0 {
        result.AverageBodySize = float64(totalBodySize) / float64(result.TotalRequests)
    }
    
    return result
}

func formatResult(result *AnalysisResult, format string, session *Session) string {
    if format == "json" {
        data, _ := json.MarshalIndent(result, "", "  ")
        return string(data)
    }
    
    // Text format
    var sb strings.Builder
    
    sb.WriteString("═══════════════════════════════════════════════════════════════\n")
    sb.WriteString("                    SESSION ANALYSIS REPORT                    \n")
    sb.WriteString("═══════════════════════════════════════════════════════════════\n\n")
    
    sb.WriteString(fmt.Sprintf("Session Name: %s\n", session.Name))
    sb.WriteString(fmt.Sprintf("Session ID: %s\n", session.ID))
    sb.WriteString(fmt.Sprintf("Created: %s\n", session.CreatedAt.Format("2006-01-02 15:04:05")))
    sb.WriteString(fmt.Sprintf("Total Requests: %d\n\n", result.TotalRequests))
    
    sb.WriteString("─ HTTP Methods ────────────────────────────────────────────────\n")
    for method, count := range result.Methods {
        sb.WriteString(fmt.Sprintf("  %-10s: %d\n", method, count))
    }
    
    sb.WriteString("\n─ Status Codes ────────────────────────────────────────────────\n")
    for code, count := range result.StatusCodes {
        sb.WriteString(fmt.Sprintf("  %-10d: %d\n", code, count))
    }
    
    sb.WriteString("\n─ Statistics ─────────────────────────────────────────────────\n")
    sb.WriteString(fmt.Sprintf("  Unique Hosts    : %d\n", result.UniqueHosts))
    sb.WriteString(fmt.Sprintf("  Average Body Size: %.2f bytes\n", result.AverageBodySize))
    
    sb.WriteString("\n─ Request Timeline ────────────────────────────────────────────\n")
    for i, req := range session.Requests {
        if i >= 20 {
            sb.WriteString(fmt.Sprintf("  ... and %d more requests\n", len(session.Requests)-20))
            break
        }
        status := ""
        if req.Response != nil {
            status = fmt.Sprintf("→ %d", req.Response.StatusCode)
        }
        sb.WriteString(fmt.Sprintf("  %2d. %s %s %s\n", i+1, req.Method, req.URL, status))
    }
    
    sb.WriteString("\n═══════════════════════════════════════════════════════════════\n")
    
    return sb.String()
}
