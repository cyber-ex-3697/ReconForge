package report

import (
    "encoding/json"
    "fmt"
    "os"
    "strings"
    "time"
)

type ReportData struct {
    Target       string    `json:"target"`
    Timestamp    time.Time `json:"timestamp"`
    Duration     string    `json:"duration"`
    Subdomains   []string  `json:"subdomains"`
    LiveHosts    []string  `json:"live_hosts"`
    URLs         []string  `json:"urls"`
    Vulnerabilities []Vulnerability `json:"vulnerabilities"`
}

type Vulnerability struct {
    URL       string `json:"url"`
    TemplateID string `json:"template_id"`
    Severity  string `json:"severity"`
    Name      string `json:"name"`
}

func Generate(data *ReportData, outputDir string) error {
    // Generate JSON
    jsonData, _ := json.MarshalIndent(data, "", "  ")
    os.WriteFile(outputDir+"/report.json", jsonData, 0644)
    
    // Generate HTML
    html := generateHTML(data)
    os.WriteFile(outputDir+"/report.html", []byte(html), 0644)
    
    return nil
}

func generateHTML(data *ReportData) string {
    return fmt.Sprintf(`<!DOCTYPE html>
<html><head><title>ReconForge Report</title></head>
<body>
<h1>Scan Report for %s</h1>
<p>Date: %s | Duration: %s</p>
<h2>Subdomains (%d)</h2>
<pre>%s</pre>
<h2>Live Hosts (%d)</h2>
<pre>%s</pre>
</body></html>`,
        data.Target, data.Timestamp.Format("2006-01-02 15:04:05"),
        data.Duration, len(data.Subdomains), strings.Join(data.Subdomains, "\n"),
        len(data.LiveHosts), strings.Join(data.LiveHosts, "\n"))
}
