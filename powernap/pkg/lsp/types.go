package lsp

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/purpose168/charm-experimental-packages-cn/powernap/pkg/lsp/protocol"
	"github.com/purpose168/charm-experimental-packages-cn/powernap/pkg/transport"
)

// OffsetEncoding 表示用于文本文档偏移的字符编码。
type OffsetEncoding int

const (
	// UTF8 编码 - 字节。
	UTF8 OffsetEncoding = iota
	// UTF16 编码 - LSP 的默认编码。
	UTF16
	// UTF32 编码 - 代码点。
	UTF32
)

// Client 表示与语言服务器的 LSP 客户端连接。
type Client struct {
	ID               string
	Name             string
	conn             *transport.Connection
	ctx              context.Context
	cancel           context.CancelFunc
	initialized      atomic.Bool
	shutdown         atomic.Bool
	capabilities     protocol.ServerCapabilities
	offsetEncoding   OffsetEncoding
	rootURI          string
	workspaceFolders []protocol.WorkspaceFolder
	config           map[string]any
	initOptions      map[string]any
}

// ClientConfig 表示创建新 LSP 客户端的配置。
type ClientConfig struct {
	Command          string
	Args             []string
	RootURI          string
	WorkspaceFolders []protocol.WorkspaceFolder
	InitOptions      map[string]any
	Settings         map[string]any
	Environment      map[string]string
	Timeout          time.Duration
}
