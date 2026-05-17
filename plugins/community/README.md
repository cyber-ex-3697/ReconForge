# Community Plugins Directory

This directory contains community-contributed plugins.

## How to Contribute

1. Fork the repository
2. Create your plugin
3. Submit a pull request

## Plugin Template

```go
package main

import (
    "context"
    "reconforge/plugins"
)

type CommunityPlugin struct {
    plugins.BasePlugin
}

func init() {
    plugin := &CommunityPlugin{
        BasePlugin: plugins.BasePlugin{
            Info: plugins.PluginInfo{
                Name:    "CommunityPlugin",
                Version: "1.0.0",
                Author:  "Your Name",
                Phase:   "subdomain",
            },
        },
    }
    RegisterPlugin(plugin)
}

Installation

cp myplugin.so plugins/community/


