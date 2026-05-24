package config

import (
    "os"
    "gopkg.in/yaml.v3"
)

type Config struct {
    Threads    int    `yaml:"threads"`
    Deep       bool   `yaml:"deep"`
    Timeout    int    `yaml:"timeout"`
    Retries    int    `yaml:"retries"`
    RateLimit  int    `yaml:"rate_limit"`
    LogLevel   string `yaml:"log_level"`
    OutputDir  string `yaml:"output_dir"`
    API        struct {
        ChaosKey   string `yaml:"chaos_key"`
        GitHubToken string `yaml:"github_token"`
        ShodanKey  string `yaml:"shodan_key"`
    } `yaml:"api"`
}

func Load(path string) (*Config, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return DefaultConfig(), nil
    }
    
    var cfg Config
    yaml.Unmarshal(data, &cfg)
    return &cfg, nil
}

func DefaultConfig() *Config {
    return &Config{
        Threads:   50,
        Deep:      false,
        Timeout:   30,
        Retries:   3,
        RateLimit: 10,
        LogLevel:  "info",
    }
}
