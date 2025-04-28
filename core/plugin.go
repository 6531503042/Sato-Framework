package core

import (
	"fmt"
	"sync"

	"github.com/gofiber/fiber/v2"
)

// PluginState represents the state of a plugin
type PluginState int

const (
	PluginStateUninitialized PluginState = iota
	PluginStateInitialized
	PluginStateStarted
	PluginStateStopped
)

// Plugin interface defines the contract for all plugins
type Plugin interface {
	// Register is called when the plugin is initialized
	Register(app *App) error
	// Start is called when the plugin is started
	Start() error
	// Stop is called when the plugin is stopped
	Stop() error
	// GetName returns the name of the plugin
	GetName() string
	// GetVersion returns the version of the plugin
	GetVersion() string
	// GetState returns the current state of the plugin
	GetState() PluginState
}

// PluginRegistry manages all registered plugins
type PluginRegistry struct {
	plugins map[string]Plugin
	mu      sync.RWMutex
}

// NewPluginRegistry creates a new plugin registry
func NewPluginRegistry() *PluginRegistry {
	return &PluginRegistry{
		plugins: make(map[string]Plugin),
	}
}

// RegisterPlugin adds a new plugin to the registry
func (r *PluginRegistry) RegisterPlugin(plugin Plugin) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.plugins[plugin.GetName()]; exists {
		return ErrPluginAlreadyRegistered
	}

	r.plugins[plugin.GetName()] = plugin
	return nil
}

// GetPlugin retrieves a plugin by name
func (r *PluginRegistry) GetPlugin(name string) (Plugin, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	plugin, exists := r.plugins[name]
	if !exists {
		return nil, ErrPluginNotFound
	}

	return plugin, nil
}

// InitializePlugins initializes all registered plugins
func (r *PluginRegistry) InitializePlugins(app *App) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, plugin := range r.plugins {
		if err := plugin.Register(app); err != nil {
			return fmt.Errorf("failed to initialize plugin %s: %v", plugin.GetName(), err)
		}
	}

	return nil
}

// StartPlugins starts all registered plugins
func (r *PluginRegistry) StartPlugins() error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, plugin := range r.plugins {
		if err := plugin.Start(); err != nil {
			return fmt.Errorf("failed to start plugin %s: %v", plugin.GetName(), err)
		}
	}

	return nil
}

// StopPlugins stops all registered plugins
func (r *PluginRegistry) StopPlugins() error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, plugin := range r.plugins {
		if err := plugin.Stop(); err != nil {
			return fmt.Errorf("failed to stop plugin %s: %v", plugin.GetName(), err)
		}
	}

	return nil
}

// Plugin errors
var (
	ErrPluginAlreadyRegistered = fiber.NewError(fiber.StatusConflict, "plugin already registered")
	ErrPluginNotFound         = fiber.NewError(fiber.StatusNotFound, "plugin not found")
) 