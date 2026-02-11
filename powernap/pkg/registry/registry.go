// Package registry 提供了一个用于管理多个语言服务器实例的注册表。
package registry

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"maps"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/purpose168/charm-experimental-packages-cn/powernap/pkg/config"
	"github.com/purpose168/charm-experimental-packages-cn/powernap/pkg/lsp"
	"github.com/purpose168/charm-experimental-packages-cn/powernap/pkg/lsp/protocol"
)

// Registry 管理多个语言服务器实例。
type Registry struct {
	mu      sync.RWMutex
	clients map[string]*lsp.Client
	configs map[string]*config.ServerConfig
	logger  *slog.Logger
}

// New 创建一个新的注册表。
func New() *Registry {
	return &Registry{
		clients: make(map[string]*lsp.Client),
		configs: make(map[string]*config.ServerConfig),
		logger:  slog.Default(),
	}
}

// NewWithLogger 创建一个带有自定义日志记录器的新注册表。
func NewWithLogger(logger *slog.Logger) *Registry {
	return &Registry{
		clients: make(map[string]*lsp.Client),
		configs: make(map[string]*config.ServerConfig),
		logger:  logger,
	}
}

// LoadConfig 从配置管理器加载服务器配置。
func (r *Registry) LoadConfig(cfg *config.Manager) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	servers := cfg.GetServers()
	maps.Copy(r.configs, servers)
	return nil
}

// StartServer 为给定的名称和项目路径启动语言服务器。
func (r *Registry) StartServer(ctx context.Context, name string, projectPath string) (*lsp.Client, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if server is already running
	if client, exists := r.clients[name]; exists {
		return client, nil
	}

	// Get server configuration
	serverCfg, exists := r.configs[name]
	if !exists {
		return nil, fmt.Errorf("no configuration found for server: %s", name)
	}

	// Find project root
	rootPath := r.findProjectRoot(projectPath, serverCfg.RootMarkers)
	if rootPath == "" {
		// Check if server supports single file mode
		if !serverCfg.SingleFileSupport {
			return nil, fmt.Errorf("language server %s requires a project root with one of: %v", name, serverCfg.RootMarkers)
		}
		rootPath = projectPath
	}

	// Create workspace folders
	workspaceFolders := []protocol.WorkspaceFolder{
		{
			URI:  "file://" + rootPath,
			Name: filepath.Base(rootPath),
		},
	}

	// Create client configuration
	clientCfg := lsp.ClientConfig{
		Command:          serverCfg.Command,
		Args:             serverCfg.Args,
		RootURI:          "file://" + rootPath,
		WorkspaceFolders: workspaceFolders,
		InitOptions:      serverCfg.InitOptions,
		Settings:         serverCfg.Settings,
		Environment:      serverCfg.Environment,
	}

	// Create and initialize client
	client, err := lsp.NewClient(clientCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	// Initialize the client
	if err := client.Initialize(ctx, serverCfg.EnableSnippets); err != nil {
		return nil, fmt.Errorf("failed to initialize client: %w", err)
	}

	// Store the client
	r.clients[name] = client

	r.logger.Info("Started language server", "name", name, "root", rootPath)
	return client, nil
}

// StopServer 停止运行中的语言服务器。
func (r *Registry) StopServer(ctx context.Context, name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	client, exists := r.clients[name]
	if !exists {
		return fmt.Errorf("server not running: %s", name)
	}

	// Shutdown the client
	if err := client.Shutdown(ctx); err != nil &&
		!errors.Is(err, io.EOF) &&
		!errors.Is(err, context.Canceled) &&
		err.Error() != "signal: killed" {
		r.logger.Error("Failed to shutdown server", "name", name, "error", err)
	}

	// Send exit notification
	if err := client.Exit(); err != nil {
		r.logger.Error("Failed to exit server", "name", name, "error", err)
	}

	// Remove from registry
	delete(r.clients, name)

	r.logger.Info("Stopped language server", "name", name)
	return nil
}

// RestartServer 重启语言服务器。
func (r *Registry) RestartServer(ctx context.Context, name string, projectPath string) (*lsp.Client, error) {
	// Stop the server if it's running
	if err := r.StopServer(ctx, name); err != nil {
		// Ignore error if server wasn't running
		r.logger.Debug("Server was not running", "name", name)
	}

	// Start the server
	return r.StartServer(ctx, name, projectPath)
}

// GetClient 通过名称返回运行中的客户端。
func (r *Registry) GetClient(name string) (*lsp.Client, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	client, exists := r.clients[name]
	return client, exists
}

// GetClientsForFile 返回给定文件的所有适用客户端。
// 这允许多个语言服务器处理同一文件类型（例如，Go 文件的 gopls 和 golangci-lint）。
func (r *Registry) GetClientsForFile(ctx context.Context, filePath string) ([]*lsp.Client, error) {
	// Convert to absolute path
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Detect language from file
	language := string(lsp.DetectLanguage(absPath))
	if language == "" {
		return nil, fmt.Errorf("unsupported file type: %s", filepath.Ext(absPath))
	}

	r.mu.RLock()

	// Find all servers that support this language
	var serverNames []string
	for name, cfg := range r.configs {
		for _, ft := range cfg.FileTypes {
			// Match by extension or language ID
			ext := filepath.Ext(absPath)
			if ft == language || ft == strings.TrimPrefix(ext, ".") || "."+ft == ext {
				serverNames = append(serverNames, name)
				break
			}
		}
	}

	r.mu.RUnlock()

	if len(serverNames) == 0 {
		return nil, fmt.Errorf("no language servers found for language: %s", language)
	}

	var clients []*lsp.Client
	projectDir := filepath.Dir(absPath)

	// Start or get each server
	for _, serverName := range serverNames {
		// Check if server is already running
		if client, exists := r.GetClient(serverName); exists {
			clients = append(clients, client)
		} else {
			// Start the server
			client, err := r.StartServer(ctx, serverName, projectDir)
			if err != nil {
				r.logger.Warn("Failed to start server", "name", serverName, "error", err)
				continue
			}
			clients = append(clients, client)
		}
	}

	if len(clients) == 0 {
		return nil, fmt.Errorf("failed to start any language servers for language: %s", language)
	}

	return clients, nil
}

// GetClientForFile 返回给定文件的单个适用客户端。
// 这是一个返回第一个可用客户端的便捷方法。
// 对于多服务器支持，请使用 GetClientsForFile 代替。
func (r *Registry) GetClientForFile(ctx context.Context, filePath string) (*lsp.Client, error) {
	clients, err := r.GetClientsForFile(ctx, filePath)
	if err != nil {
		return nil, err
	}

	if len(clients) == 0 {
		return nil, fmt.Errorf("no clients available")
	}

	return clients[0], nil
}

// ListClients 返回运行中客户端的列表。
func (r *Registry) ListClients() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.clients))
	for name := range r.clients {
		names = append(names, name)
	}

	return names
}

// StopAll 停止所有运行中的语言服务器。
func (r *Registry) StopAll(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	var errs []error

	for name, client := range r.clients {
		if err := client.Shutdown(ctx); err != nil {
			errs = append(errs, fmt.Errorf("failed to shutdown %s: %w", name, err))
		}

		if err := client.Exit(); err != nil {
			errs = append(errs, fmt.Errorf("failed to exit %s: %w", name, err))
		}
	}

	// Clear all clients
	r.clients = make(map[string]*lsp.Client)

	if len(errs) > 0 {
		return fmt.Errorf("errors stopping servers: %v", errs)
	}

	return nil
}

// findProjectRoot 根据根标记查找项目根目录。
func (r *Registry) findProjectRoot(startPath string, rootMarkers []string) string {
	currentPath := startPath

	for {
		// Check for root markers
		for _, marker := range rootMarkers {
			markerPath := filepath.Join(currentPath, marker)
			if fileExists(markerPath) {
				return currentPath
			}
		}

		// Move up one directory
		parentPath := filepath.Dir(currentPath)
		if parentPath == currentPath {
			// Reached filesystem root
			break
		}

		currentPath = parentPath
	}

	return ""
}

// fileExists 检查文件是否存在。
func fileExists(path string) bool {
	// Use os.Stat to check if the file/directory exists
	if _, err := os.Stat(path); err != nil {
		return false
	}
	return true
}
