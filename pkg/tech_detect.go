package livehost

import (
    "encoding/json"
    "os/exec"
    "strings"
)

type Technology struct {
    Name     string   `json:"name"`
    Version  string   `json:"version"`
    Categories []string `json:"categories"`
}

type TechDetector struct {
    threads int
}

func NewTechDetector(threads int) *TechDetector {
    return &TechDetector{
        threads: threads,
    }
}

func (t *TechDetector) Detect(url string) ([]Technology, error) {
    cmd := exec.Command("httpx", "-u", url, "-silent", "-tech-detect", "-json")
    output, err := cmd.Output()
    if err != nil {
        return nil, err
    }
    
    var result struct {
        Technologies []Technology `json:"tech"`
    }
    
    if err := json.Unmarshal(output, &result); err != nil {
        return nil, err
    }
    
    return result.Technologies, nil
}

func (t *TechDetector) DetectBatch(urls []string) (map[string][]Technology, error) {
    results := make(map[string][]Technology)
    
    input := strings.Join(urls, "\n")
    cmd := exec.Command("httpx", "-silent", "-tech-detect", "-json", "-threads", string(rune(t.threads)))
    cmd.Stdin = strings.NewReader(input)
    output, err := cmd.Output()
    if err != nil {
        return nil, err
    }
    
    for _, line := range strings.Split(string(output), "\n") {
        if line == "" {
            continue
        }
        var techResult struct {
            URL          string       `json:"url"`
            Technologies []Technology `json:"tech"`
        }
        if err := json.Unmarshal([]byte(line), &techResult); err == nil {
            results[techResult.URL] = techResult.Technologies
        }
    }
    return results, nil
}
