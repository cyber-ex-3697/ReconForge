package config

func GetDefaultWordlist() string {
    wordlists := []string{
        "/usr/share/wordlists/subdomains.txt",
        "/usr/share/seclists/Discovery/DNS/subdomains-top1million-5000.txt",
        "/usr/share/wordlists/dirb/common.txt",
    }
    
    for _, wl := range wordlists {
        if fileExists(wl) {
            return wl
        }
    }
    return wordlists[0]
}

func fileExists(path string) bool {
    _, err := os.Stat(path)
    return err == nil
}
