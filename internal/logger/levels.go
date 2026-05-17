package logger

import (
    "strings"
)

func ParseLevel(s string) Level {
    switch strings.ToUpper(s) {
    case "DEBUG":
        return LevelDebug
    case "INFO":
        return LevelInfo
    case "WARN":
        return LevelWarn
    case "ERROR":
        return LevelError
    default:
        return LevelInfo
    }
}

func (l Level) String() string {
    return string(l)
}

func (l Level) ShouldLog(currentLevel Level) bool {
    levels := map[Level]int{
        LevelDebug: 0,
        LevelInfo:  1,
        LevelWarn:  2,
        LevelError: 3,
    }
    return levels[l] >= levels[currentLevel]
}
