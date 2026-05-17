package config

import (
    "os"
    "gopkg.in/yaml.v3"
)

type Config struct {
    Version string     `yaml:"version"`
    Scan    ScanConfig `yaml:"scan"`
    Rate    RateConfig `yaml:"rate"`
    API     APIConfig  `yaml:"api"`
    Tools   ToolsConfig `yaml:"tools"`
}

type ScanConfig struct {
    Profile   string `yaml:"profile"`
    Threads   int    `yaml:"threads"`
    DeepMode  bool   `yaml:"deep_mode"`
    Timeout   int    `yaml:"timeout"`
    Retries   int    `yaml:"retries"`
    RateLimit int    `yaml:"rate_limit"`
}

type RateConfig struct {
    Default    int `yaml:"default"`
    Aggressive int `yaml:"aggressive"`
    Stealth    int `yaml:"stealth"`
}

type APIConfig struct {
    ChaosKey         string `yaml:"chaos_key"`
    GitHubToken      string `yaml:"github_token"`
    ShodanKey        string `yaml:"shodan_key"`
    CensysKey        string `yaml:"censys_key"`
    CensysSecret     string `yaml:"censys_secret"`
    SecurityTrailsKey string `yaml:"securitytrails_key"`
    VirusTotalKey    string `yaml:"virustotal_key"`
}

type ToolsConfig struct {
    SubfinderVersion string   `yaml:"subfinder_version"`
    NucleiVersion    string   `yaml:"nuclei_version"`
    HttpxVersion     string   `yaml:"httpx_version"`
    Wordlists        []string `yaml:"wordlists"`
    Resolvers        []string `yaml:"resolvers"`
}

func Load(path string) (*Config, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return DefaultConfig(), nil
    }

    var cfg Config
    if err := yaml.Unmarshal(data, &cfg); err != nil {
        return nil, err
    }

    return &cfg, nil
}

func DefaultConfig() *Config {
    return &Config{
        Version: "4.0.0",
        Scan: ScanConfig{
            Profile:   "standard",
            Threads:   50,
            DeepMode:  false,
            Timeout:   30,
            Retries:   3,
            RateLimit: 10,
        },
        Rate: RateConfig{
            Default:    50,
            Aggressive: 100,
            Stealth:    10,
        },
        API: APIConfig{},
        Tools: ToolsConfig{
            SubfinderVersion: "2.14.0",
            NucleiVersion:    "3.1.0",
            HttpxVersion:     "1.3.8",
            Wordlists:        []string{"/usr/share/wordlists/subdomains.txt"},
            Resolvers:        []string{"1.1.1.1", "8.8.8.8"},
        },
    }
}
