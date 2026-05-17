package subdomain

import (
    "bufio"
    "os"
    "os/exec"
    "strings"
)

type BruteForcer struct {
    target   string
    wordlist string
    resolvers []string
}

func NewBruteForcer(target string, wordlist string) *BruteForcer {
    return &BruteForcer{
        target:   target,
        wordlist: wordlist,
        resolvers: []string{"1.1.1.1", "8.8.8.8"},
    }
}

func (b *BruteForcer) SetResolvers(resolvers []string) {
    b.resolvers = resolvers
}

func (b *BruteForcer) RunShuffledns() ([]string, error) {
    cmd := exec.Command("shuffledns", "-d", b.target, "-w", b.wordlist, "-silent")
    output, err := cmd.Output()
    if err != nil {
        return nil, err
    }
    
    var results []string
    for _, s := range strings.Split(string(output), "\n") {
        if s != "" {
            results = append(results, s)
        }
    }
    return results, nil
}

func (b *BruteForcer) RunPuredns() ([]string, error) {
    cmd := exec.Command("puredns", "resolve", b.wordlist, "-d", b.target, "-r", "resolvers.txt", "-w", "/tmp/resolved.txt")
    if err := cmd.Run(); err != nil {
        return nil, err
    }
    
    file, err := os.Open("/tmp/resolved.txt")
    if err != nil {
        return nil, err
    }
    defer file.Close()
    
    var results []string
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        results = append(results, scanner.Text())
    }
    return results, nil
}

func (b *BruteForcer) Run() ([]string, error) {
    return b.RunShuffledns()
}
