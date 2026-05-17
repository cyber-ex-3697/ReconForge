# ReconForge - Plugin Development Guide

## Overview

Plugins extend ReconForge functionality. They can be written in Go and loaded dynamically.

## Plugin Interface

```go
type Plugin interface {
    GetInfo() PluginInfo
    Init(config map[string]interface{}) error
    Run(ctx context.Context, target string, input interface{}) (*PluginResult, error)
    Cleanup() error
}

Creating a Plugin

Step 1: Create Plugin File

package main

import (
    "context"
    "reconforge/plugins"
)

type MyPlugin struct {
    plugins.BasePlugin
}

func NewMyPlugin() *MyPlugin {
    return &MyPlugin{
        BasePlugin: plugins.BasePlugin{
            Info: plugins.PluginInfo{
                Name:        "MyPlugin",
                Version:     "1.0.0",
                Author:      "Your Name",
                Description: "Plugin description",
                Phase:       "subdomain",
            },
        },
    }
}

func (p *MyPlugin) Run(ctx context.Context, target string, input interface{}) (*plugins.PluginResult, error) {
    // Your logic here
    
    return &plugins.PluginResult{
        Success: true,
        Data:    result,
    }, nil
}

func (p *MyPlugin) Cleanup() error {
    return nil
}

Step 2: Build Plugin

go build -buildmode=plugin -o myplugin.so myplugin.go

Step 3: Install Plugin

cp myplugin.so plugins/community/

Plugin Phases

Phase		When Executed

subdomain	After subdomain enumeration
livehost	After live host detection
url		After URL discovery
vulnerability	After vulnerability scan
report		During report generation

Example: Custom Subdomain Plugin

package main

import (
    "context"
    "fmt"
    "reconforge/plugins"
)

type CustomSubdomainPlugin struct {
    plugins.BasePlugin
}

func (p *CustomSubdomainPlugin) Run(ctx context.Context, target string, input interface{}) (*plugins.PluginResult, error) {
    // Custom subdomain logic
    subdomains := []string{
        "api." + target,
        "admin." + target,
        "cdn." + target,
    }
    
    return &plugins.PluginResult{
        Success: true,
        Data:    subdomains,
    }, nil
}

Example: Notification Plugin

package main

import (
    "context"
    "fmt"
    "reconforge/plugins"
)

type NotifyPlugin struct {
    plugins.BasePlugin
}

func (p *NotifyPlugin) Run(ctx context.Context, target string, input interface{}) (*plugins.PluginResult, error) {
    vulnerabilities := input.([]Vulnerability)
    
    if len(vulnerabilities) > 0 {
        fmt.Printf("Found %d vulnerabilities on %s\n", len(vulnerabilities), target)
    }
    
    return &plugins.PluginResult{
        Success: true,
    }, nil
}


Testing Plugins

# Test plugin
./reconforge --plugin-test myplugin.so

# Run with plugin
./reconforge -t example.com --plugin myplugin.so


Plugin Configuration

Plugins can access config:

func (p *MyPlugin) Init(config map[string]interface{}) error {
    apiKey := config["api_key"].(string)
    return nil
}

Best Practices


Keep plugins focused - One task per plugin

Handle errors gracefully - Don't crash main program

Use output directory - Save files to p.OutputDir

Clean up resources - Implement Cleanup()

Document your plugin - Include README



Submitting Plugins


Fork repository

Add plugin to plugins/community/

Submit pull request

Include documentation


Community Plugins


Takeover Plugin

Slack Notifier

Discord Webhook
