package main

import (
    "flag"
    "fmt"
    "net/http"
    "os"
    "os/exec"
    "path/filepath"
)

func main() {
    var reportPath string
    var port int
    var openBrowser bool
    
    flag.StringVar(&reportPath, "report", "", "Path to HTML report file")
    flag.IntVar(&port, "port", 8080, "HTTP server port")
    flag.BoolVar(&openBrowser, "open", true, "Open browser automatically")
    flag.Parse()
    
    fmt.Println("=== ReconForge Report Viewer ===\n")
    
    if reportPath == "" {
        fmt.Println("Usage: report_viewer -report <report.html>")
        fmt.Println("\nOptions:")
        fmt.Println("  -report   Path to HTML report file")
        fmt.Println("  -port     HTTP server port (default: 8080)")
        fmt.Println("  -open     Open browser automatically (default: true)")
        os.Exit(1)
    }
    
    // Check if report exists
    if _, err := os.Stat(reportPath); os.IsNotExist(err) {
        fmt.Printf("[!] Report file not found: %s\n", reportPath)
        os.Exit(1)
    }
    
    // Get absolute path
    absPath, err := filepath.Abs(reportPath)
    if err != nil {
        fmt.Printf("[!] Error getting absolute path: %v\n", err)
        os.Exit(1)
    }
    
    reportDir := filepath.Dir(absPath)
    reportFile := filepath.Base(absPath)
    
    fmt.Printf("[✓] Report: %s\n", absPath)
    fmt.Printf("[*] Serving directory: %s\n", reportDir)
    
    // Start HTTP server
    http.Handle("/", http.FileServer(http.Dir(reportDir)))
    
    url := fmt.Sprintf("http://localhost:%d/%s", port, reportFile)
    
    fmt.Printf("[*] Starting server on port %d\n", port)
    
    if openBrowser {
        fmt.Printf("[*] Opening browser at %s\n", url)
        openBrowserURL(url)
    }
    
    fmt.Printf("\n[✓] Report available at: %s\n", url)
    fmt.Println("Press Ctrl+C to stop")
    
    http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func openBrowserURL(url string) {
    var err error
    switch os.Getenv("OSTYPE") {
    case "linux-gnu":
        err = exec.Command("xdg-open", url).Start()
    case "darwin":
        err = exec.Command("open", url).Start()
    default:
        err = exec.Command("xdg-open", url).Start()
    }
    if err != nil {
        fmt.Printf("[!] Could not open browser: %v\n", err)
        // Try fallback
        exec.Command("firefox", url).Start()
    }
}
