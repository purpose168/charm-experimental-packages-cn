package lsp

import (
	"encoding/json"
)

// Message 表示一个 JSON-RPC 2.0 消息。
type Message struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      int32           `json:"id,omitempty"`
	Method  string          `json:"method,omitempty"`
	Params  json.RawMessage `json:"params,omitempty"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *ResponseError  `json:"error,omitempty"`
}

// ResponseError 表示一个 JSON-RPC 2.0 错误。
type ResponseError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// NewRequest 创建一个新的 JSON-RPC 2.0 请求消息。
func NewRequest(id int32, method string, params any) (*Message, error) {
	paramsJSON, err := json.Marshal(params)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	return &Message{
		JSONRPC: "2.0",
		ID:      id,
		Method:  method,
		Params:  paramsJSON,
	}, nil
}

// NewNotification 创建一个新的 JSON-RPC 2.0 通知消息。
func NewNotification(method string, params any) (*Message, error) {
	paramsJSON, err := json.Marshal(params)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	return &Message{
		JSONRPC: "2.0",
		Method:  method,
		Params:  paramsJSON,
	}, nil
}
