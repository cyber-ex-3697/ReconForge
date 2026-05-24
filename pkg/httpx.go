package livehost

import (
    "encoding/json"
    "os/exec"
    "strings"
)

type HTTPXResult struct {
    URL         string `json:"url"`
    StatusCode  int    `json:"status_code"`
    Title       string `json:"title"`
    ContentType string `json:"content_type"`
    WebServer   string `json:"webserver"`
    ContentLength int  `json:"content_length"`
    ResponseTime  int64 `json:"response_time"`
    Technologies []string `json:"technologies"`
}

type HTTPXWrapper struct {
    threads int
    timeout int
}

func NewHTTPXWrapper(threads int) *HTTPXWrapper {
    return &HTTPXWrapper{
        threads: threads,
        timeout: 10,
    }
}

func (h *HTTPXWrapper) Scan(domains []string) ([]HTTPXResult, error) {
    input := strings.Join(domains, "\n")
    cmd := exec.Command("httpx", "-silent", "-json", "-status-code", "-title", "-tech-detect", "-threads", string(rune(h.threads)))
    cmd.Stdin = strings.NewReader(input)
    output, err := cmd.Output()
    if err != nil {
        return nil, err
    }
    
    var results []HTTPXResult
    for _, line := range strings.Split(string(output), "\n") {
        if line == "" {
            continue
        }
        var result HTTPXResult
        if err := json.Unmarshal([]byte(line), &result); err == nil {
            results = append(results, result)
        }
    }
    return results, nil
}

func (h *HTTPXWrapper) GetLiveHosts(domains []string) ([]string, error) {
    results, err := h.Scan(domains)
    if err != nil {
        return nil, err
    }
    
    var live []string
    for _, r := range results {
        if r.StatusCode == 200 || r.StatusCode == 301 || r.StatusCode == 302 {
            live = append(live, r.URL)
        }
    }
    return live, nil
}
