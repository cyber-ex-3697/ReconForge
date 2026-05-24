package proxy

import (
    "net/http"
    "net/url"
    "sync"
    "time"
)

type ProxyRotator struct {
    proxies   []string
    current   int
    mu        sync.Mutex
    lastUsed  time.Time
}

func NewProxyRotator() *ProxyRotator {
    return &ProxyRotator{
        proxies: []string{
            "http://proxy1:8080",
            "http://proxy2:8080",
            "socks5://proxy3:1080",
        },
        current:  0,
        lastUsed: time.Now(),
    }
}

func (p *ProxyRotator) GetClient() *http.Client {
    p.mu.Lock()
    defer p.mu.Unlock()
    
    // Rotate proxy every 10 requests
    if time.Since(p.lastUsed) > 30*time.Second {
        p.current = (p.current + 1) % len(p.proxies)
        p.lastUsed = time.Now()
    }
    
    proxyURL, _ := url.Parse(p.proxies[p.current])
    return &http.Client{
        Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)},
        Timeout:   30 * time.Second,
    }
}
