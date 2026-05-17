package report

import (
    "encoding/json"
    "os"
    "strings"
)

type SearchableIndex struct {
    Vulnerabilities []Vulnerability `json:"vulnerabilities"`
    LiveHosts       []string        `json:"live_hosts"`
    Subdomains      []string        `json:"subdomains"`
}

type SearchEngine struct {
    indexFile string
    data      *SearchableIndex
}

func NewSearchEngine(indexFile string) *SearchEngine {
    return &SearchEngine{
        indexFile: indexFile,
        data:      &SearchableIndex{},
    }
}

func (s *SearchEngine) BuildIndex(data *ReportData) error {
    s.data = &SearchableIndex{
        Vulnerabilities: data.Vulnerabilities,
        LiveHosts:       data.LiveHosts,
        Subdomains:      data.Subdomains,
    }
    
    jsonData, err := json.MarshalIndent(s.data, "", "  ")
    if err != nil {
        return err
    }
    return os.WriteFile(s.indexFile, jsonData, 0644)
}

func (s *SearchEngine) Search(query string) *SearchableIndex {
    results := &SearchableIndex{}
    queryLower := strings.ToLower(query)
    
    for _, v := range s.data.Vulnerabilities {
        if strings.Contains(strings.ToLower(v.Name), queryLower) ||
           strings.Contains(strings.ToLower(v.TemplateID), queryLower) {
            results.Vulnerabilities = append(results.Vulnerabilities, v)
        }
    }
    
    for _, host := range s.data.LiveHosts {
        if strings.Contains(strings.ToLower(host), queryLower) {
            results.LiveHosts = append(results.LiveHosts, host)
        }
    }
    
    for _, sub := range s.data.Subdomains {
        if strings.Contains(strings.ToLower(sub), queryLower) {
            results.Subdomains = append(results.Subdomains, sub)
        }
    }
    
    return results
}

func (s *SearchEngine) LoadIndex() error {
    data, err := os.ReadFile(s.indexFile)
    if err != nil {
        return err
    }
    return json.Unmarshal(data, s.data)
}
