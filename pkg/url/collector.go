package url

import (
    "bufio"
    "os"
    "os/exec"
    "strings"
)

type Collector struct {
    target string
}

type URLResult struct {
    URLs []string
    Total int
}

func NewCollector(target string) *Collector {
    return &Collector{
        target: target,
    }
}

func (c *Collector) Collect() (*URLResult, error) {
    urls := make(map[string]bool)
    
    // GAU
    cmd := exec.Command("gau", "--subs", c.target)
    output, err := cmd.Output()
    if err == nil {
        for _, u := range strings.Split(string(output), "\n") {
            if u != "" {
                urls[u] = true
            }
        }
    }
    
    // Waybackurls
    cmd = exec.Command("waybackurls", c.target)
    output, err = cmd.Output()
    if err == nil {
        for _, u := range strings.Split(string(output), "\n") {
            if u != "" {
                urls[u] = true
            }
        }
    }
    
    result := &URLResult{
        URLs:  make([]string, 0, len(urls)),
        Total: len(urls),
    }
    for u := range urls {
        result.URLs = append(result.URLs, u)
    }
    return result, nil
}

func (c *Collector) SaveToFile(filename string, urls []string) error {
    content := strings.Join(urls, "\n")
    return os.WriteFile(filename, []byte(content), 0644)
}

func (c *Collector) LoadFromFile(filename string) ([]string, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()
    
    var urls []string
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if line != "" {
            urls = append(urls, line)
        }
    }
    return urls, scanner.Err()
}
