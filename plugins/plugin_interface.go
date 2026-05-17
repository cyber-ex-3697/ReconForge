package plugins

import (
    "context"
)

// PluginInfo contains metadata about a plugin
type PluginInfo struct {
    Name        string   `json:"name"`
    Version     string   `json:"version"`
    Author      string   `json:"author"`
    Description string   `json:"description"`
    Phase       string   `json:"phase"`        // subdomain, livehost, url, vulnerability, etc.
    Dependencies []string `json:"dependencies"` // Required tools
}

// PluginResult contains the output from a plugin
type PluginResult struct {
    Success     bool        `json:"success"`
    Data        interface{} `json:"data"`
    Error       string      `json:"error,omitempty"`
    OutputFile  string      `json:"output_file,omitempty"`
}

// Plugin is the interface that all plugins must implement
type Plugin interface {
    // GetInfo returns metadata about the plugin
    GetInfo() PluginInfo
    
    // Init initializes the plugin with configuration
    Init(config map[string]interface{}) error
    
    // Run executes the plugin
    Run(ctx context.Context, target string, input interface{}) (*PluginResult, error)
    
    // Cleanup performs cleanup operations
    Cleanup() error
}

// BasePlugin provides a base implementation for plugins
type BasePlugin struct {
    Info     PluginInfo
    Config   map[string]interface{}
    OutputDir string
}

func (b *BasePlugin) GetInfo() PluginInfo {
    return b.Info
}

func (b *BasePlugin) Init(config map[string]interface{}) error {
    b.Config = config
    return nil
}

func (b *BasePlugin) Cleanup() error {
    return nil
}

// PluginManager manages all loaded plugins
type PluginManager struct {
    plugins map[string]Plugin
}

func NewPluginManager() *PluginManager {
    return &PluginManager{
        plugins: make(map[string]Plugin),
    }
}

func (pm *PluginManager) Register(name string, plugin Plugin) {
    pm.plugins[name] = plugin
}

func (pm *PluginManager) Get(name string) (Plugin, bool) {
    plugin, ok := pm.plugins[name]
    return plugin, ok
}

func (pm *PluginManager) GetAll() map[string]Plugin {
    return pm.plugins
}

func (pm *PluginManager) GetByPhase(phase string) []Plugin {
    var result []Plugin
    for _, plugin := range pm.plugins {
        if plugin.GetInfo().Phase == phase {
            result = append(result, plugin)
        }
    }
    return result
}
