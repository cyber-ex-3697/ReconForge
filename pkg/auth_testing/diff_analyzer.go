package auth_testing

import (
    "crypto/md5"
    "encoding/hex"
    "hash/fnv"
    "math"
    "regexp"
    "strings"
)

// DiffAnalyzer performs semantic response diffing
type DiffAnalyzer struct {
    shingleSize int
    hashSize    int
    ignorePatterns []*regexp.Regexp
}

// DiffResult contains comparison results
type DiffResult struct {
    Similarity    float64            `json:"similarity"`
    IsDifferent   bool               `json:"is_different"`
    Differences   []string           `json:"differences"`
    LengthDiff    int                `json:"length_diff"`
    HashDiff      bool               `json:"hash_diff"`
}

// NewDiffAnalyzer creates a new diff analyzer
func NewDiffAnalyzer() *DiffAnalyzer {
    return &DiffAnalyzer{
        shingleSize: 9,
        hashSize:    100,
        ignorePatterns: []*regexp.Regexp{
            regexp.MustCompile(`"timestamp":"[^"]+"`),
            regexp.MustCompile(`"date":"[^"]+"`),
            regexp.MustCompile(`"nonce":"[^"]+"`),
            regexp.MustCompile(`"session":"[^"]+"`),
            regexp.MustCompile(`"csrf":"[^"]+"`),
            regexp.MustCompile(`"token":"[^"]+"`),
        },
    }
}

// Compare compares two responses
func (da *DiffAnalyzer) Compare(response1, response2 string) *DiffResult {
    // Normalize both responses
    norm1 := da.normalize(response1)
    norm2 := da.normalize(response2)
    
    result := &DiffResult{
        Similarity:  0.0,
        IsDifferent: false,
        Differences: make([]string, 0),
        LengthDiff:  len(norm1) - len(norm2),
        HashDiff:    da.hash(norm1) != da.hash(norm2),
    }
    
    // MinHash similarity
    similarity := da.minHashSimilarity(norm1, norm2)
    result.Similarity = similarity
    
    // Determine if meaningfully different
    result.IsDifferent = similarity < 0.85
    
    // Extract differences
    if result.IsDifferent {
        result.Differences = da.extractDifferences(norm1, norm2)
    }
    
    return result
}

// normalize normalizes response for comparison
func (da *DiffAnalyzer) normalize(response string) string {
    // Remove whitespace
    response = strings.TrimSpace(response)
    
    // Apply ignore patterns
    for _, pattern := range da.ignorePatterns {
        response = pattern.ReplaceAllString(response, `"$1":"[REDACTED]"`)
    }
    
    // Normalize numbers
    numberRegex := regexp.MustCompile(`:\s*\d+`)
    response = numberRegex.ReplaceAllString(response, `: [NUMBER]`)
    
    // Normalize IDs
    idRegex := regexp.MustCompile(`"id":"[^"]+"`)
    response = idRegex.ReplaceAllString(response, `"id":"[ID]"`)
    
    // Normalize emails
    emailRegex := regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
    response = emailRegex.ReplaceAllString(response, "[EMAIL]")
    
    return response
}

// hash computes hash of response
func (da *DiffAnalyzer) hash(response string) string {
    hash := md5.Sum([]byte(response))
    return hex.EncodeToString(hash[:])
}

// minHashSimilarity computes MinHash similarity between two strings
func (da *DiffAnalyzer) minHashSimilarity(s1, s2 string) float64 {
    shingles1 := da.getShingles(s1)
    shingles2 := da.getShingles(s2)
    
    if len(shingles1) == 0 && len(shingles2) == 0 {
        return 1.0
    }
    
    if len(shingles1) == 0 || len(shingles2) == 0 {
        return 0.0
    }
    
    // Compute Jaccard similarity
    set1 := make(map[uint64]bool)
    set2 := make(map[uint64]bool)
    
    for _, h := range shingles1 {
        set1[h] = true
    }
    for _, h := range shingles2 {
        set2[h] = true
    }
    
    intersection := 0
    for h := range set1 {
        if set2[h] {
            intersection++
        }
    }
    
    union := len(set1) + len(set2) - intersection
    
    return float64(intersection) / float64(union)
}

// getShingles extracts shingles from string
func (da *DiffAnalyzer) getShingles(s string) []uint64 {
    if len(s) < da.shingleSize {
        return []uint64{da.hashToUint64(s)}
    }
    
    shingles := make([]uint64, 0)
    for i := 0; i <= len(s)-da.shingleSize; i++ {
        shingle := s[i : i+da.shingleSize]
        shingles = append(shingles, da.hashToUint64(shingle))
    }
    
    return shingles
}

// hashToUint64 converts string to uint64 hash
func (da *DiffAnalyzer) hashToUint64(s string) uint64 {
    h := fnv.New64a()
    h.Write([]byte(s))
    return h.Sum64()
}

// extractDifferences extracts meaningful differences
func (da *DiffAnalyzer) extractDifferences(s1, s2 string) []string {
    differences := make([]string, 0)
    
    // Simple line-by-line comparison
    lines1 := strings.Split(s1, "\n")
    lines2 := strings.Split(s2, "\n")
    
    maxLen := len(lines1)
    if len(lines2) > maxLen {
        maxLen = len(lines2)
    }
    
    for i := 0; i < maxLen; i++ {
        var line1, line2 string
        if i < len(lines1) {
            line1 = lines1[i]
        }
        if i < len(lines2) {
            line2 = lines2[i]
        }
        
        if line1 != line2 {
            if line1 != "" && line2 != "" {
                differences = append(differences, fmt.Sprintf("Line %d differs", i+1))
            } else if line1 != "" {
                differences = append(differences, fmt.Sprintf("Line %d only in first response", i+1))
            } else {
                differences = append(differences, fmt.Sprintf("Line %d only in second response", i+1))
            }
        }
        
        if len(differences) >= 10 {
            differences = append(differences, "... more differences")
            break
        }
    }
    
    return differences
}

// CompareWithMinHash performs MinHash-based comparison
func (da *DiffAnalyzer) CompareWithMinHash(response1, response2 string) float64 {
    return da.minHashSimilarity(response1, response2)
}

// GetStructuralDiff finds structural differences between JSON responses
func (da *DiffAnalyzer) GetStructuralDiff(json1, json2 string) []string {
    differences := make([]string, 0)
    
    // Simple JSON key comparison
    keyRegex := regexp.MustCompile(`"([a-zA-Z0-9_]+)"\s*:`)
    
    keys1 := make(map[string]bool)
    keys2 := make(map[string]bool)
    
    for _, match := range keyRegex.FindAllStringSubmatch(json1, -1) {
        if len(match) > 1 {
            keys1[match[1]] = true
        }
    }
    
    for _, match := range keyRegex.FindAllStringSubmatch(json2, -1) {
        if len(match) > 1 {
            keys2[match[1]] = true
        }
    }
    
    for key := range keys1 {
        if !keys2[key] {
            differences = append(differences, fmt.Sprintf("Key '%s' only in first response", key))
        }
    }
    
    for key := range keys2 {
        if !keys1[key] {
            differences = append(differences, fmt.Sprintf("Key '%s' only in second response", key))
        }
    }
    
    return differences
}
