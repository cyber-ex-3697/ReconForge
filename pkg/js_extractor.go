package url

import (
    "io"
    "net/http"
    "regexp"
    "strings"
)

type JSExtractor struct {
    patterns []string
}

func NewJSExtractor() *JSExtractor {
    return &JSExtractor{
        patterns: []string{
            `https?://[^\s"']+\.js`,
            `/[a-zA-Z0-9/_\-]+\.js`,
            `api\.js`,
            `config\.js`,
            `endpoints\.js`,
            `/js/`,
            `/assets/js/`,
            `/static/js/`,
        },
    }
}

func (j *JSExtractor) ExtractFromURL(url string) ([]string, error) {
    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }
    
    var jsFiles []string
    for _, pattern := range j.patterns {
        re := regexp.MustCompile(pattern)
        matches := re.FindAllString(string(body), -1)
        jsFiles = append(jsFiles, matches...)
    }
    return jsFiles, nil
}

func (j *JSExtractor) ExtractFromHTML(html string) []string {
    var jsFiles []string
    for _, pattern := range j.patterns {
        re := regexp.MustCompile(pattern)
        matches := re.FindAllString(html, -1)
        jsFiles = append(jsFiles, matches...)
    }
    return jsFiles
}

func (j *JSExtractor) AddPattern(pattern string) {
    j.patterns = append(j.patterns, pattern)
}

func (j *JSExtractor) ExtractEndpointsFromJS(jsURL string) ([]string, error) {
    resp, err := http.Get(jsURL)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }
    
    // Find API endpoints in JS
    apiPatterns := []string{
        `/api/[a-zA-Z0-9/_\-?=&]+`,
        `/v1/[a-zA-Z0-9/_\-?=&]+`,
        `/v2/[a-zA-Z0-9/_\-?=&]+`,
        `/graphql`,
        `/rest/[a-zA-Z0-9/_\-?=&]+`,
    }
    
    var endpoints []string
    for _, pattern := range apiPatterns {
        re := regexp.MustCompile(pattern)
        matches := re.FindAllString(string(body), -1)
        endpoints = append(endpoints, matches...)
    }
    return endpoints, nil
}
