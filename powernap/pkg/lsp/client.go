// Package lsp 提供了语言服务器协议（Language Server Protocol，LSP）的客户端实现。
package lsp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/purpose168/charm-experimental-packages-cn/powernap/pkg/lsp/protocol"
	"github.com/purpose168/charm-experimental-packages-cn/powernap/pkg/transport"
)

// LSP 方法常量。
const (
	MethodInitialize                         = "initialize"
	MethodInitialized                        = "initialized"
	MethodShutdown                           = "shutdown"
	MethodExit                               = "exit"
	MethodTextDocumentDidOpen                = "textDocument/didOpen"
	MethodTextDocumentDidChange              = "textDocument/didChange"
	MethodTextDocumentDidSave                = "textDocument/didSave"
	MethodTextDocumentDidClose               = "textDocument/didClose"
	MethodTextDocumentCompletion             = "textDocument/completion"
	MethodTextDocumentHover                  = "textDocument/hover"
	MethodTextDocumentDefinition             = "textDocument/definition"
	MethodTextDocumentReferences             = "textDocument/references"
	MethodTextDocumentDiagnostic             = "textDocument/publishDiagnostics"
	MethodWorkspaceConfiguration             = "workspace/configuration"
	MethodWorkspaceDidChangeConfiguration    = "workspace/didChangeConfiguration"
	MethodWorkspaceDidChangeWorkspaceFolders = "workspace/didChangeWorkspaceFolders"
	MethodWorkspaceDidChangeWatchedFiles     = "workspace/didChangeWatchedFiles"
)

// NewClient 使用给定的配置创建一个新的 LSP 客户端。
func NewClient(config ClientConfig) (*Client, error) {
	ctx, cancel := context.WithCancel(context.Background())

	client := &Client{
		ID:               config.Command, // Will be updated after initialization
		Name:             config.Command,
		ctx:              ctx,
		cancel:           cancel,
		rootURI:          config.RootURI,
		workspaceFolders: config.WorkspaceFolders,
		config:           config.Settings,
		initOptions:      config.InitOptions,
		offsetEncoding:   UTF16, // Default to UTF16
	}

	// Start the language server process
	stream, err := startServerProcess(ctx, config)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to start language server: %w", err)
	}

	// Create transport connection
	conn, err := transport.NewConnection(ctx, stream, slog.Default())
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to create connection: %w", err)
	}

	client.conn = conn

	// Register handlers for server-initiated requests
	client.setupHandlers()

	return client, nil
}

// Kill 取消连接上下文。
func (c *Client) Kill() { c.cancel() }

// Initialize 向语言服务器发送初始化请求。
func (c *Client) Initialize(ctx context.Context, enableSnippets bool) error {
	if c.initialized.Load() {
		return fmt.Errorf("client already initialized")
	}

	// Extract root path from URI
	rootPath := ""
	if c.rootURI != "" {
		rootPath = strings.TrimPrefix(c.rootURI, "file://")
	}

	// Prepare workspace folders - some servers don't like nil
	workspaceFolders := c.workspaceFolders
	if workspaceFolders == nil {
		workspaceFolders = []protocol.WorkspaceFolder{}
	}

	initParams := map[string]any{
		"processId": os.Getpid(),
		"clientInfo": map[string]any{
			"name":    "powernap",
			"version": "0.1.0",
		},
		"locale":                "en-us",
		"rootPath":              rootPath, // Deprecated but some servers still use it
		"rootUri":               c.rootURI,
		"capabilities":          c.makeClientCapabilities(enableSnippets),
		"workspaceFolders":      workspaceFolders,
		"initializationOptions": c.initOptions, // Use the client's init options
		"trace":                 "off",         // Can be "off", "messages", or "verbose"
	}

	// Log the initialization params for debugging
	paramsJSON, _ := json.MarshalIndent(initParams, "", "  ")
	slog.Debug("Sending initialize request", "params", string(paramsJSON))

	var result protocol.InitializeResult
	err := c.conn.Call(ctx, MethodInitialize, initParams, &result)
	if err != nil {
		return fmt.Errorf("initialize request failed: %w", err)
	}

	// Store server capabilities
	c.capabilities = result.Capabilities

	// Handle offset encoding
	if result.OffsetEncoding != "" {
		switch result.OffsetEncoding {
		case "utf-8":
			c.offsetEncoding = UTF8
		case "utf-16":
			c.offsetEncoding = UTF16
		case "utf-32":
			c.offsetEncoding = UTF32
		}
	}

	// Send initialized notification
	err = c.conn.Notify(ctx, MethodInitialized, map[string]any{})
	if err != nil {
		return fmt.Errorf("initialized notification failed: %w", err)
	}

	c.initialized.Store(true)

	// For gopls, send workspace/didChangeConfiguration to ensure it's ready
	// This helps gopls properly set up its workspace views
	if strings.Contains(c.Name, "gopls") {
		configParams := map[string]any{
			"settings": c.config,
		}
		_ = c.conn.Notify(ctx, MethodWorkspaceDidChangeConfiguration, configParams)

		// Also send workspace/didChangeWatchedFiles to trigger gopls to scan the workspace
		// This helps with the "no views" error
		if c.rootURI != "" {
			changesParams := map[string]any{
				"changes": []map[string]any{
					{
						"uri":  c.rootURI,
						"type": 1, // Created
					},
				},
			}
			_ = c.conn.Notify(ctx, "workspace/didChangeWatchedFiles", changesParams)
		}
	}

	return nil
}

// Shutdown 向语言服务器发送关闭请求。
func (c *Client) Shutdown(ctx context.Context) error {
	if c.shutdown.Load() {
		return nil
	}

	err := c.conn.Call(ctx, MethodShutdown, nil, nil)
	if err != nil {
		return fmt.Errorf("shutdown request failed: %w", err)
	}

	c.shutdown.Store(true)
	return nil
}

// Exit 向语言服务器发送退出通知。
func (c *Client) Exit() error {
	err := c.conn.Notify(c.ctx, MethodExit, nil)
	if err != nil {
		return fmt.Errorf("exit notification failed: %w", err)
	}

	c.cancel()
	return nil
}

// GetCapabilities 返回服务器能力。
func (c *Client) GetCapabilities() protocol.ServerCapabilities {
	return c.capabilities
}

// IsInitialized 返回客户端是否已初始化。
func (c *Client) IsInitialized() bool {
	return c.initialized.Load()
}

// IsRunning 返回客户端连接是否仍然活跃。
func (c *Client) IsRunning() bool {
	return c.conn != nil && c.conn.IsConnected() && c.initialized.Load() && !c.shutdown.Load()
}

// RegisterNotificationHandler 注册一个处理服务器发起的通知的处理器。
func (c *Client) RegisterNotificationHandler(method string, handler transport.NotificationHandler) {
	if c.conn != nil {
		c.conn.RegisterNotificationHandler(method, handler)
	}
}

// RegisterHandler 注册一个处理服务器发起的请求的处理器。
func (c *Client) RegisterHandler(method string, handler transport.Handler) {
	if c.conn != nil {
		c.conn.RegisterHandler(method, handler)
	}
}

// NotifyDidOpenTextDocument 通知服务器文档已打开。
func (c *Client) NotifyDidOpenTextDocument(ctx context.Context, uri string, languageID string, version int, text string) error {
	if !c.initialized.Load() {
		return fmt.Errorf("client not initialized")
	}

	params := protocol.DidOpenTextDocumentParams{
		TextDocument: protocol.TextDocumentItem{
			URI:        protocol.DocumentURI(uri),
			LanguageID: protocol.LanguageKind(languageID),
			Version:    int32(version), //nolint:gosec
			Text:       text,
		},
	}

	// Log what we're sending for debugging
	slog.Debug("Sending textDocument/didOpen",
		"uri", uri,
		"languageId", languageID,
		"version", version,
		"textLength", len(text))

	return c.conn.Notify(ctx, MethodTextDocumentDidOpen, params) //nolint:wrapcheck
}

// NotifyDidCloseTextDocument 通知服务器文档已关闭。
func (c *Client) NotifyDidCloseTextDocument(ctx context.Context, uri string) error {
	if !c.initialized.Load() {
		return fmt.Errorf("client not initialized")
	}

	params := protocol.DidCloseTextDocumentParams{
		TextDocument: protocol.TextDocumentIdentifier{
			URI: protocol.DocumentURI(uri),
		},
	}

	return c.conn.Notify(ctx, MethodTextDocumentDidClose, params) //nolint:wrapcheck
}

// NotifyDidChangeTextDocument 通知服务器文档已更改。
func (c *Client) NotifyDidChangeTextDocument(ctx context.Context, uri string, version int, changes []protocol.TextDocumentContentChangeEvent) error {
	if !c.initialized.Load() {
		return fmt.Errorf("client not initialized")
	}

	params := protocol.DidChangeTextDocumentParams{
		TextDocument: protocol.VersionedTextDocumentIdentifier{
			Version: int32(version), //nolint:gosec
			TextDocumentIdentifier: protocol.TextDocumentIdentifier{
				URI: protocol.DocumentURI(uri),
			},
		},
		ContentChanges: changes,
	}

	return c.conn.Notify(ctx, MethodTextDocumentDidChange, params) //nolint:wrapcheck
}

// NotifyDidChangeWatchedFiles 通知服务器被监视的文件已更改。
func (c *Client) NotifyDidChangeWatchedFiles(ctx context.Context, changes []protocol.FileEvent) error {
	if !c.initialized.Load() {
		return fmt.Errorf("client not initialized")
	}

	params := protocol.DidChangeWatchedFilesParams{
		Changes: changes,
	}

	return c.conn.Notify(ctx, MethodWorkspaceDidChangeWatchedFiles, params) //nolint:wrapcheck
}

// NotifyWorkspaceDidChangeConfiguration 通知服务器工作区配置已更改。
func (c *Client) NotifyWorkspaceDidChangeConfiguration(ctx context.Context, settings any) error {
	if !c.initialized.Load() {
		return fmt.Errorf("client not initialized")
	}

	params := map[string]any{
		"settings": settings,
	}

	return c.conn.Notify(ctx, MethodWorkspaceDidChangeConfiguration, params) //nolint:wrapcheck
}

// RequestCompletion 请求给定位置的补全项。
func (c *Client) RequestCompletion(ctx context.Context, uri string, position protocol.Position) (*protocol.CompletionList, error) {
	if !c.initialized.Load() {
		return nil, fmt.Errorf("client not initialized")
	}

	params := protocol.CompletionParams{
		Context: protocol.CompletionContext{
			TriggerKind: protocol.Invoked,
		},
		TextDocumentPositionParams: protocol.TextDocumentPositionParams{
			TextDocument: protocol.TextDocumentIdentifier{
				URI: protocol.DocumentURI(uri),
			},
			Position: position,
		},
	}

	var result any
	err := c.conn.Call(ctx, MethodTextDocumentCompletion, params, &result)
	if err != nil {
		return nil, fmt.Errorf("completion request failed: %w", err)
	}

	// Parse the result - can be CompletionList or []CompletionItem
	var completionList protocol.CompletionList

	switch v := result.(type) {
	case map[string]any:
		// It's a CompletionList
		data, err := json.Marshal(v)
		if err != nil {
			return nil, err //nolint:wrapcheck
		}
		if err := json.Unmarshal(data, &completionList); err != nil {
			return nil, err //nolint:wrapcheck
		}
	case []any:
		// It's an array of CompletionItem
		data, err := json.Marshal(v)
		if err != nil {
			return nil, err //nolint:wrapcheck
		}
		var items []protocol.CompletionItem
		if err := json.Unmarshal(data, &items); err != nil {
			return nil, err //nolint:wrapcheck
		}
		completionList.Items = items
		completionList.IsIncomplete = false
	}

	return &completionList, nil
}

// RequestHover 请求给定位置的悬停信息。
func (c *Client) RequestHover(ctx context.Context, uri string, position protocol.Position) (*protocol.Hover, error) {
	if !c.initialized.Load() {
		return nil, fmt.Errorf("client not initialized")
	}

	params := map[string]any{
		"textDocument": map[string]any{
			"uri": uri,
		},
		"position": position,
	}

	var result protocol.Hover
	err := c.conn.Call(ctx, MethodTextDocumentHover, params, &result)
	if err != nil {
		return nil, fmt.Errorf("hover request failed: %w", err)
	}

	return &result, nil
}

// FindReferences 查找给定位置符号的所有引用。
func (c *Client) FindReferences(ctx context.Context, filepath string, line, character int, includeDeclaration bool) ([]protocol.Location, error) {
	uri := string(protocol.URIFromPath(filepath))
	params := protocol.ReferenceParams{
		TextDocumentPositionParams: protocol.TextDocumentPositionParams{
			TextDocument: protocol.TextDocumentIdentifier{
				URI: protocol.DocumentURI(uri),
			},
			Position: protocol.Position{
				Line:      uint32(line),      //nolint:gosec
				Character: uint32(character), //nolint:gosec
			},
		},
		Context: protocol.ReferenceContext{
			IncludeDeclaration: includeDeclaration,
		},
	}

	var result []protocol.Location
	err := c.conn.Call(ctx, MethodTextDocumentReferences, params, &result)
	if err != nil {
		return nil, fmt.Errorf("find references request failed: %w", err)
	}
	return result, nil
}

// setupHandlers 注册处理服务器发起的请求的处理器。
func (c *Client) setupHandlers() {
	// 处理 workspace/configuration 请求
	c.conn.RegisterHandler(MethodWorkspaceConfiguration, func(_ context.Context, _ string, params json.RawMessage) (any, error) {
		var configParams protocol.ConfigurationParams
		if err := json.Unmarshal(params, &configParams); err != nil {
			return nil, err //nolint:wrapcheck
		}

		// Return configuration for each requested item
		result := make([]any, len(configParams.Items))
		for i := range configParams.Items {
			result[i] = c.config
		}

		return result, nil
	})

	// 处理其他常见的服务器请求
	// 根据需要添加更多处理器
}

// makeClientCapabilities 创建初始化所需的客户端能力。
func (c *Client) makeClientCapabilities(enableSnippets bool) map[string]any {
	return map[string]any{
		"textDocument": map[string]any{
			"synchronization": map[string]any{
				"dynamicRegistration": true,
				"willSave":            true,
				"willSaveWaitUntil":   true,
				"didSave":             true,
			},
			"completion": map[string]any{
				"dynamicRegistration": true,
				"completionItem": map[string]any{
					"snippetSupport":          enableSnippets,
					"commitCharactersSupport": true,
					"documentationFormat":     []string{"markdown", "plaintext"},
					"deprecatedSupport":       true,
					"preselectSupport":        true,
					"insertReplaceSupport":    true,
					"tagSupport": map[string]any{
						"valueSet": []int{1}, // Deprecated
					},
					"resolveSupport": map[string]any{
						"properties": []string{"documentation", "detail", "additionalTextEdits"},
					},
				},
				"contextSupport": true,
			},
			"hover": map[string]any{
				"dynamicRegistration": true,
				"contentFormat":       []string{"markdown", "plaintext"},
			},
			"definition": map[string]any{
				"dynamicRegistration": true,
				"linkSupport":         true,
			},
			"references": map[string]any{
				"dynamicRegistration": true,
			},
			"documentHighlight": map[string]any{
				"dynamicRegistration": true,
			},
			"documentSymbol": map[string]any{
				"dynamicRegistration":               true,
				"hierarchicalDocumentSymbolSupport": true,
			},
			"formatting": map[string]any{
				"dynamicRegistration": true,
			},
			"rangeFormatting": map[string]any{
				"dynamicRegistration": true,
			},
			"rename": map[string]any{
				"dynamicRegistration": true,
				"prepareSupport":      true,
			},
			"publishDiagnostics": map[string]any{
				"relatedInformation":     true,
				"versionSupport":         true,
				"tagSupport":             map[string]any{"valueSet": []int{1, 2}},
				"codeDescriptionSupport": true,
				"dataSupport":            true,
			},
			"codeAction": map[string]any{
				"dynamicRegistration": true,
				"codeActionLiteralSupport": map[string]any{
					"codeActionKind": map[string]any{
						"valueSet": []string{
							"quickfix",
							"refactor",
							"refactor.extract",
							"refactor.inline",
							"refactor.rewrite",
							"source",
							"source.organizeImports",
						},
					},
				},
				"isPreferredSupport": true,
				"dataSupport":        true,
				"resolveSupport": map[string]any{
					"properties": []string{"edit"},
				},
			},
		},
		"workspace": map[string]any{
			"applyEdit": true,
			"workspaceEdit": map[string]any{
				"documentChanges":       true,
				"resourceOperations":    []string{"create", "rename", "delete"},
				"failureHandling":       "textOnlyTransactional",
				"normalizesLineEndings": true,
			},
			"didChangeConfiguration": map[string]any{
				"dynamicRegistration": true,
			},
			"didChangeWatchedFiles": map[string]any{
				"dynamicRegistration":    true,
				"relativePatternSupport": true,
			},
			"symbol": map[string]any{
				"dynamicRegistration": true,
			},
			"configuration":    true,
			"workspaceFolders": true,
			"fileOperations": map[string]any{
				"dynamicRegistration": true,
				"didCreate":           true,
				"willCreate":          true,
				"didRename":           true,
				"willRename":          true,
				"didDelete":           true,
				"willDelete":          true,
			},
		},
		"window": map[string]any{
			"workDoneProgress": true,
			"showMessage": map[string]any{
				"messageActionItem": map[string]any{
					"additionalPropertiesSupport": true,
				},
			},
			"showDocument": map[string]any{
				"support": true,
			},
		},
		"general": map[string]any{
			"regularExpressions": map[string]any{
				"engine":  "ECMAScript",
				"version": "ES2020",
			},
			"markdown": map[string]any{
				"parser":  "marked",
				"version": "1.1.0",
			},
			"positionEncodings": []string{"utf-16"},
		},
	}
}

// startServerProcess 启动语言服务器进程。
func startServerProcess(ctx context.Context, config ClientConfig) (io.ReadWriteCloser, error) {
	cmd := exec.CommandContext(ctx, config.Command, config.Args...) //nolint:gosec

	// 设置环境变量
	if config.Environment != nil {
		cmd.Env = os.Environ()
		for k, v := range config.Environment {
			cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
		}
	}

	// 创建管道
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdin pipe: %w", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	// 创建 stderr 管道以捕获错误消息
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	// 启动进程
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start process: %w", err)
	}

	// 监控 stderr
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := stderr.Read(buf)
			if err != nil {
				if err != io.EOF {
					slog.Error("Error reading stderr", "error", err)
				}
				break
			}
			if n > 0 {
				slog.Error("Language server stderr", "command", config.Command, "output", string(buf[:n]))
			}
		}
	}()

	closer := &processCloser{
		cmd:    cmd,
		stdin:  stdin,
		stdout: stdout,
		stderr: stderr,
	}

	return transport.NewStreamTransport(stdout, stdin, closer), nil
}

type processCloser struct {
	cmd       *exec.Cmd
	stdin     io.WriteCloser
	stdout    io.ReadCloser
	stderr    io.ReadCloser
	closeOnce sync.Once
	closeErr  error
}

func (c *processCloser) Close() error {
	c.closeOnce.Do(func() {
		errs := []error{
			c.stdin.Close(),
			c.stdout.Close(),
			c.stderr.Close(),
		}

		done := make(chan error, 1)
		go func() {
			done <- c.cmd.Wait()
		}()

		timeout := time.After(5 * time.Second)
		select {
		case err := <-done:
			errs = append(errs, err)
		case <-timeout:
			errs = append(errs, c.cmd.Process.Kill())
			<-done
		}

		c.closeErr = errors.Join(errs...)
	})
	return c.closeErr
}
