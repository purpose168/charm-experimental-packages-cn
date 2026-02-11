// Package config 表示语言服务器的配置管理。
package config

import (
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/mitchellh/mapstructure"
)

//go:embed lsps.json
var lspsJSON []byte

// ServerConfig 表示语言服务器的配置。
type ServerConfig struct {
	Command           string            `mapstructure:"command" json:"command"`
	Args              []string          `mapstructure:"args" json:"args,omitempty"`
	FileTypes         []string          `mapstructure:"filetypes" json:"filetypes"`
	RootMarkers       []string          `mapstructure:"root_markers" json:"root_markers"`
	Environment       map[string]string `mapstructure:"environment" json:"environment,omitempty"`
	Settings          map[string]any    `mapstructure:"settings" json:"settings,omitempty"`
	InitOptions       map[string]any    `mapstructure:"init_options" json:"init_options,omitempty"`
	EnableSnippets    bool              `mapstructure:"enable_snippets" json:"-"`
	SingleFileSupport bool              `mapstructure:"single_file_support" json:"-"`
}

// Config 表示整体配置。
type Config struct {
	Servers map[string]*ServerConfig `mapstructure:"servers"`
}

// Manager 管理配置加载和访问。
type Manager struct {
	config *Config
}

// NewManager 创建一个新的配置管理器。
func NewManager() *Manager {
	return &Manager{
		config: &Config{
			Servers: make(map[string]*ServerConfig),
		},
	}
}

// LoadDefaults 从嵌入的 JSON 加载默认服务器配置。
func (m *Manager) LoadDefaults() error {
	servers := make(map[string]*ServerConfig)
	if err := json.Unmarshal(lspsJSON, &servers); err != nil {
		return fmt.Errorf("failed to parse embedded lsps.json: %w", err)
	}

	m.config.Servers = servers
	m.applyDefaults()
	return nil
}

// GetServers 返回所有服务器配置。
func (m *Manager) GetServers() map[string]*ServerConfig {
	return m.config.Servers
}

// GetServer 返回特定的服务器配置。
func (m *Manager) GetServer(name string) (*ServerConfig, bool) {
	server, exists := m.config.Servers[name]
	return server, exists
}

// AddServer 添加或更新服务器配置。
func (m *Manager) AddServer(name string, config *ServerConfig) {
	m.config.Servers[name] = config
}

// RemoveServer 移除服务器配置。
func (m *Manager) RemoveServer(name string) {
	delete(m.config.Servers, name)
}

// applyDefaults 应用默认值到服务器配置。
func (m *Manager) applyDefaults() {
	for name, server := range m.config.Servers {
		if server.RootMarkers == nil {
			server.RootMarkers = []string{".git"}
		}

		if server.Environment == nil {
			server.Environment = make(map[string]string)
		}

		if server.Settings == nil {
			server.Settings = make(map[string]any)
		}
		_, server.EnableSnippets = snippetSupport[name]
		_, server.SingleFileSupport = singleFileSupport[name]
	}
}

// LoadFromMap 从映射加载配置（对测试有用）。
func (m *Manager) LoadFromMap(data map[string]any) error {
	var config Config
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:  &config,
		TagName: "mapstructure",
	})
	if err != nil {
		return fmt.Errorf("failed to create decoder: %w", err)
	}

	if err := decoder.Decode(data); err != nil {
		return fmt.Errorf("failed to decode config: %w", err)
	}

	m.config = &config
	m.applyDefaults()

	return nil
}
