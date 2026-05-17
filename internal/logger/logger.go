package logger

import (
    "fmt"
    "time"
)

type Level string

const (
    LevelDebug Level = "DEBUG"
    LevelInfo  Level = "INFO"
    LevelWarn  Level = "WARN"
    LevelError Level = "ERROR"
)

type Logger struct {
    debug   bool
    jsonLog bool
}

func New(debug bool) *Logger {
    return &Logger{
        debug:   debug,
        jsonLog: false,
    }
}

func (l *Logger) log(level Level, msg string) {
    timestamp := time.Now().Format("15:04:05")
    
    if l.jsonLog {
        fmt.Printf(`{"time":"%s","level":"%s","msg":"%s"}`, timestamp, level, msg)
    } else {
        switch level {
        case LevelDebug:
            fmt.Printf("[DEBUG] %s %s\n", timestamp, msg)
        case LevelInfo:
            fmt.Printf("[INFO] %s %s\n", timestamp, msg)
        case LevelWarn:
            fmt.Printf("[WARN] %s %s\n", timestamp, msg)
        case LevelError:
            fmt.Printf("[ERROR] %s %s\n", timestamp, msg)
        }
    }
}

func (l *Logger) Debug(msg string) {
    if l.debug {
        l.log(LevelDebug, msg)
    }
}

func (l *Logger) Info(msg string) {
    l.log(LevelInfo, msg)
}

func (l *Logger) Warn(msg string) {
    l.log(LevelWarn, msg)
}

func (l *Logger) Error(msg string) {
    l.log(LevelError, msg)
}

func (l *Logger) Success(msg string) {
    fmt.Printf("[✓] %s\n", msg)
}
