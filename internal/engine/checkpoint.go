package engine

import (
    "encoding/json"
    "os"
    "time"
)

type Checkpoint struct {
    Target         string    `json:"target"`
    Phase          Phase     `json:"phase"`
    Subdomains     []string  `json:"subdomains"`
    LiveHosts      []string  `json:"live_hosts"`
    URLs           []string  `json:"urls"`
    Vulnerabilities []string `json:"vulnerabilities"`
    Timestamp      time.Time `json:"timestamp"`
}

func (s *ScanState) SaveCheckpoint(filename string, result *ScanResult) error {
    checkpoint := &Checkpoint{
        Target:         result.Target,
        Phase:          s.GetPhase(),
        Subdomains:     result.Subdomains,
        LiveHosts:      result.LiveHosts,
        URLs:           result.URLs,
        Vulnerabilities: result.Vulnerabilities,
        Timestamp:      time.Now(),
    }
    
    data, err := json.MarshalIndent(checkpoint, "", "  ")
    if err != nil {
        return err
    }
    
    return os.WriteFile(filename, data, 0644)
}

func LoadCheckpoint(filename string) (*Checkpoint, error) {
    data, err := os.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    
    var checkpoint Checkpoint
    if err := json.Unmarshal(data, &checkpoint); err != nil {
        return nil, err
    }
    
    return &checkpoint, nil
}
