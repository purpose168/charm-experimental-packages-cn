// Package transport 提供用于 LSP 通信的 JSON-RPC 2.0 传输层。
package transport

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"sync"
	"sync/atomic"

	"github.com/sourcegraph/jsonrpc2"
)

// Connection 表示到语言服务器的托管连接。
type Connection struct {
	conn      jsonrpc2.JSONRPC2
	transport *Transport
	router    *Router
	logger    *slog.Logger

	// 状态管理
	closed   atomic.Bool
	closeMu  sync.Mutex
	closeErr error

	// 请求跟踪
	requestMu sync.Mutex
	requests  map[jsonrpc2.ID]chan *Message
}

// NewConnection 创建一个新的托管连接。
func NewConnection(ctx context.Context, stream io.ReadWriteCloser, logger *slog.Logger) (*Connection, error) {
	c := &Connection{
		router:   NewRouter(),
		logger:   logger,
		requests: make(map[jsonrpc2.ID]chan *Message),
	}

	// 创建 JSON-RPC 连接
	conn := jsonrpc2.NewConn(
		ctx,
		jsonrpc2.NewBufferedStream(stream, jsonrpc2.VSCodeObjectCodec{}),
		jsonrpc2.HandlerWithError(c.handleRequest),
	)

	c.conn = conn
	c.transport = NewWithConn(conn)

	return c, nil
}

// Call 向语言服务器发出请求并等待响应。
func (c *Connection) Call(ctx context.Context, method string, params any, result any) error {
	if c.closed.Load() {
		return fmt.Errorf("connection is closed")
	}

	return c.conn.Call(ctx, method, params, result) //nolint:wrapcheck
}

// Notify 向语言服务器发送通知。
func (c *Connection) Notify(ctx context.Context, method string, params any) error {
	if c.closed.Load() {
		return fmt.Errorf("connection is closed")
	}

	return c.conn.Notify(ctx, method, params) //nolint:wrapcheck
}

// handleRequest 处理来自语言服务器的传入请求。
func (c *Connection) handleRequest(ctx context.Context, _ *jsonrpc2.Conn, req *jsonrpc2.Request) (any, error) {
	if c.logger != nil {
		c.logger.Debug("Handling request", "method", req.Method)
	}

	return c.router.Route(ctx, req)
}

// RegisterHandler 为特定方法注册一个处理程序。
func (c *Connection) RegisterHandler(method string, handler Handler) {
	c.router.Handle(method, handler)
}

// RegisterNotificationHandler 注册一个通知处理程序。
func (c *Connection) RegisterNotificationHandler(method string, handler NotificationHandler) {
	c.router.HandleNotification(method, handler)
}

// Close 关闭连接。
func (c *Connection) Close() error {
	c.closeMu.Lock()
	defer c.closeMu.Unlock()

	if c.closed.Load() {
		return c.closeErr
	}

	c.closed.Store(true)

	// 关闭 JSON-RPC 连接
	if c.conn != nil {
		c.closeErr = c.conn.Close()
	}

	// 关闭任何待处理的请求
	c.requestMu.Lock()
	for _, ch := range c.requests {
		close(ch)
	}
	c.requests = nil
	c.requestMu.Unlock()

	return c.closeErr
}

// IsConnected 如果连接仍然活跃，则返回 true。
func (c *Connection) IsConnected() bool {
	return !c.closed.Load()
}
