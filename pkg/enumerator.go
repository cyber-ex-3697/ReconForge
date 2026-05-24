package subdomain

import (
    "bufio"
    "fmt"
    "os"
    "os/exec"
    "strings"
)

type Enumerator struct {
    target   string
    threads  int
    wordlist string
}

type SubdomainResult struct {
    Subdomains []string
    Total      int
}

func NewEnumerator(target string, threads int) *Enumerator {
    return &Enumerator{
        target:  target,
        threads: threads,
    }
}

func (e *Enumerator) SetWordlist(wordlist string) {
    e.wordlist = wordlist
}

func (e *Enumerator) Run() (*SubdomainResult, error) {
    result := &SubdomainResult{
        Subdomains: make([]string, 0),
    }
    
    subdomains := make(map[string]bool)
    
    // Run all enumeration methods
    e.runPassive(subdomains)
    
    if e.wordlist != "" {
        e.runBruteForce(subdomains)
    }
    
    // Convert map to slice
    for s := range subdomains {
        result.Subdomains = append(result.Subdomains, s)
    }
    result.Total = len(result.Subdomains)
    
    return result, nil
}

func (e *Enumerator) runPassive(subdomains map[string]bool) {
    // Subfinder
    cmd := exec.Command("subfinder", "-d", e.target, "-silent")
    output, _ := cmd.Output()
    for _, s := range strings.Split(string(output), "\n") {
        if s != "" && strings.Contains(s, e.target) {
            subdomains[s] = true
        }
    }
    
    // Assetfinder
    cmd = exec.Command("assetfinder", "--subs-only", e.target)
    output, _ = cmd.Output()
    for _, s := range strings.Split(string(output), "\n") {
        if s != "" && strings.Contains(s, e.target) {
            subdomains[s] = true
        }
    }
    
    // Findomain
    cmd = exec.Command("findomain", "-t", e.target, "-q")
    output, _ = cmd.Output()
    for _, s := range strings.Split(string(output), "\n") {
        if s != "" && strings.Contains(s, e.target) {
            subdomains[s] = true
        }
    }
}

func (e *Enumerator) runBruteForce(subdomains map[string]bool) {
    cmd := exec.Command("shuffledns", "-d", e.target, "-w", e.wordlist, "-silent")
    output, _ := cmd.Output()
    for _, s := range strings.Split(string(output), "\n") {
        if s != "" {
            subdomains[s] = true
        }
    }
}

func (e *Enumerator) SaveToFile(filename string, subdomains []string) error {
    content := strings.Join(subdomains, "\n")
    return os.WriteFile(filename, []byte(content), 0644)
}

func (e *Enumerator) LoadFromFile(filename string) ([]string, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()
    
    var subdomains []string
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if line != "" {
            subdomains = append(subdomains, line)
        }
    }
    return subdomains, scanner.Err()
}
