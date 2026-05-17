package config

import (
    "fmt"
    "strings"
)

func (c *Config) Validate() error {
    var errors []string

    if c.Scan.Threads < 1 || c.Scan.Threads > 500 {
        errors = append(errors, "threads must be between 1 and 500")
    }
    if c.Scan.Timeout < 1 || c.Scan.Timeout > 300 {
        errors = append(errors, "timeout must be between 1 and 300 seconds")
    }
    if c.Scan.RateLimit < 1 || c.Scan.RateLimit > 100 {
        errors = append(errors, "rate_limit must be between 1 and 100")
    }

    if len(errors) > 0 {
        return fmt.Errorf("config validation failed: %s", strings.Join(errors, "; "))
    }
    return nil
}
