package vulnerability

import (
    "math"
)

type CVSSScore struct {
    BaseScore   float64
    TemporalScore float64
    EnvironmentalScore float64
    Severity    string
}

func CalculateCVSS(av, ac, pr, ui, s, c, i, a string) *CVSSScore {
    score := &CVSSScore{}
    
    // Base score calculation (simplified)
    baseScore := 0.0
    
    // Attack Vector
    avMap := map[string]float64{"N": 0.85, "A": 0.62, "L": 0.55, "P": 0.2}
    if val, ok := avMap[av]; ok {
        baseScore += val
    }
    
    // Attack Complexity
    acMap := map[string]float64{"L": 0.77, "H": 0.44}
    if val, ok := acMap[ac]; ok {
        baseScore += val
    }
    
    // Privileges Required
    prMap := map[string]float64{"N": 0.85, "L": 0.62, "H": 0.27}
    if val, ok := prMap[pr]; ok {
        baseScore += val
    }
    
    // User Interaction
    uiMap := map[string]float64{"N": 0.85, "R": 0.62}
    if val, ok := uiMap[ui]; ok {
        baseScore += val
    }
    
    score.BaseScore = math.Min(baseScore*10, 10)
    
    // Determine severity
    if score.BaseScore >= 9.0 {
        score.Severity = "CRITICAL"
    } else if score.BaseScore >= 7.0 {
        score.Severity = "HIGH"
    } else if score.BaseScore >= 4.0 {
        score.Severity = "MEDIUM"
    } else {
        score.Severity = "LOW"
    }
    
    return score
}

func GetSeverityFromCVSS(score float64) string {
    if score >= 9.0 {
        return "CRITICAL"
    } else if score >= 7.0 {
        return "HIGH"
    } else if score >= 4.0 {
        return "MEDIUM"
    } else {
        return "LOW"
    }
}
