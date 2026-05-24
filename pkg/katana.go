package url

import (
    "os/exec"
    "strings"
)

type KatanaWrapper struct {
    depth   int
    threads int
}

func NewKatanaWrapper(depth, threads int) *KatanaWrapper {
    return &KatanaWrapper{
        depth:   depth,
        threads: threads,
    }
}

func (k *KatanaWrapper) Crawl(url string) ([]string, error) {
    cmd := exec.Command("katana", "-u", url, "-d", string(rune(k.depth)), "-silent", "-jc")
    output, err := cmd.Output()
    if err != nil {
        return nil, err
    }
    
    var urls []string
    for _, u := range strings.Split(string(output), "\n") {
        if u != "" {
            urls = append(urls, u)
        }
    }
    return urls, nil
}

func (k *KatanaWrapper) CrawlBatch(urls []string) ([]string, error) {
    input := strings.Join(urls, "\n")
    cmd := exec.Command("katana", "-list", input, "-d", string(rune(k.depth)), "-silent", "-jc", "-c", string(rune(k.threads)))
    cmd.Stdin = strings.NewReader(input)
    output, err := cmd.Output()
    if err != nil {
        return nil, err
    }
    
    var results []string
    for _, u := range strings.Split(string(output), "\n") {
        if u != "" {
            results = append(results, u)
        }
    }
    return results, nil
}

func (k *KatanaWrapper) SetDepth(depth int) {
    k.depth = depth
}

func (k *KatanaWrapper) SetThreads(threads int) {
    k.threads = threads
}
