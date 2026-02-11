package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/sourcegraph/jsonrpc2"
)

// Handler 是处理传入消息的函数。
type Handler func(ctx context.Context, method string, params json.RawMessage) (any, error)

// NotificationHandler 是处理传入通知的函数。
type NotificationHandler func(ctx context.Context, method string, params json.RawMessage)

// Router 将传入的消息路由到相应的处理程序。
type Router struct {
	mu                   sync.RWMutex
	handlers             map[string]Handler
	notificationHandlers map[string]NotificationHandler
	defaultHandler       Handler
}

// NewRouter 创建一个新的消息路由器。
func NewRouter() *Router {
	return &Router{
		handlers:             make(map[string]Handler),
		notificationHandlers: make(map[string]NotificationHandler),
	}
}

// Handle 为特定方法注册一个处理程序。
func (r *Router) Handle(method string, handler Handler) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.handlers[method] = handler
}

// HandleNotification 为特定方法注册一个通知处理程序。
func (r *Router) HandleNotification(method string, handler NotificationHandler) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.notificationHandlers[method] = handler
}

// SetDefaultHandler 为未注册的方法设置默认处理程序。
func (r *Router) SetDefaultHandler(handler Handler) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.defaultHandler = handler
}

// Route 将消息路由到相应的处理程序。
func (r *Router) Route(ctx context.Context, req *jsonrpc2.Request) (any, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// 检查是否为通知（无 ID）
	if req.ID == (jsonrpc2.ID{}) {
		if handler, ok := r.notificationHandlers[req.Method]; ok {
			handler(ctx, req.Method, *req.Params)
		}
		return nil, nil
	}

	// 是请求
	if handler, ok := r.handlers[req.Method]; ok {
		return handler(ctx, req.Method, *req.Params)
	}

	// 如果有默认处理程序，则使用
	if r.defaultHandler != nil {
		return r.defaultHandler(ctx, req.Method, *req.Params)
	}

	return nil, fmt.Errorf("no handler for method: %s", req.Method)
}

// MessageType 表示 JSON-RPC 消息的类型。
type MessageType int

const (
	// RequestMessage 是一个需要响应的请求。
	RequestMessage MessageType = iota
	// NotificationMessage 是一个不需要响应的通知。
	NotificationMessage
	// ResponseMessage 是对请求的响应。
	ResponseMessage
	// ErrorMessage 是错误响应。
	ErrorMessage
)

// ParseMessageType 确定 JSON-RPC 消息的类型。
func ParseMessageType(msg *Message) MessageType {
	if msg.Error != nil {
		return ErrorMessage
	}
	if msg.Result != nil {
		return ResponseMessage
	}
	if msg.ID != nil {
		return RequestMessage
	}
	return NotificationMessage
}
