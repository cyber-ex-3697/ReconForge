# ReconForge Plugin Development Guide

## Overview

ReconForge supports custom plugins that can extend its functionality.

## Creating a Plugin

### Step 1: Create a new Go file

Create a file like `myplugin.go` with this content:

```go
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
    // Your plugin logic here
    return &plugins.PluginResult{
        Success: true,
        Data:    "Result",
    }, nil
}


Step 2: Build the plugin

go build -buildmode=plugin -o myplugin.so myplugin.go

Step 3: Place in plugins/community/

cp myplugin.so plugins/community/

Available Phases

Phase	                 Description

subdomain              	Subdomain enumeration
livehost		Live host detection
url			URL discovery
vulnerability		Vulnerability scanning

Example Plugin

See custom_module.go for a complete working example.


