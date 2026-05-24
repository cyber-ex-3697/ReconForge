package livehost

import (
    "bufio"
    "os/exec"
    "strings"
)

type WAFInfo struct {
    URL     string
    WAFName string
    Detected bool
}

type WAFDetector struct {
    timeout int
}

func NewWAFDetector(timeout int) *WAFDetector {
    return &WAFDetector{
        timeout: timeout,
    }
}

func (w *WAFDetector) Detect(url string) (*WAFInfo, error) {
    cmd := exec.Command("wafw00f", url, "-a")
    output, err := cmd.Output()
    if err != nil {
        return &WAFInfo{
            URL:      url,
            Detected: false,
        }, nil
    }
    
    result := &WAFInfo{
        URL:      url,
        Detected: false,
    }
    
    outputStr := string(output)
    if strings.Contains(outputStr, "detected") {
        result.Detected = true
        // Extract WAF name
        scanner := bufio.NewScanner(strings.NewReader(outputStr))
        for scanner.Scan() {
            line := scanner.Text()
            if strings.Contains(line, "detects") {
                parts := strings.Fields(line)
                if len(parts) >= 2 {
                    result.WAFName = parts[len(parts)-1]
                    break
                }
            }
        }
    }
    return result, nil
}

func (w *WAFDetector) DetectBatch(urls []string) ([]*WAFInfo, error) {
    var results []*WAFInfo
    for _, url := range urls {
        info, err := w.Detect(url)
        if err == nil {
            results = append(results, info)
        }
    }
    return results, nil
}
