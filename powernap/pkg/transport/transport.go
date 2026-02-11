package transport

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"strconv"
	"strings"
	"sync"

	"github.com/sourcegraph/jsonrpc2"
)

// Transport 处理与语言服务器的低级通信。
type Transport struct {
	conn   jsonrpc2.JSONRPC2
	reader io.Reader
	writer io.Writer
	logger *slog.Logger
	mu     sync.Mutex
}

// Message 表示一个 JSON-RPC 消息。
type Message struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      *jsonrpc2.ID    `json:"id,omitempty"`
	Method  string          `json:"method,omitempty"`
	Params  json.RawMessage `json:"params,omitempty"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *jsonrpc2.Error `json:"error,omitempty"`
}

// New 创建一个新的传输。
func New(reader io.Reader, writer io.Writer, logger *slog.Logger) *Transport {
	return &Transport{
		reader: reader,
		writer: writer,
		logger: logger,
	}
}

// NewWithConn 使用现有的 JSON-RPC 连接创建一个新的传输。
func NewWithConn(conn jsonrpc2.JSONRPC2) *Transport {
	return &Transport{
		conn: conn,
	}
}

// Send 向语言服务器发送消息。
func (t *Transport) Send(ctx context.Context, msg *Message) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.conn != nil {
		// Use existing connection
		if msg.ID != nil {
			// It's a request
			var result json.RawMessage
			err := t.conn.Call(ctx, msg.Method, msg.Params, &result)
			if err != nil {
				return err //nolint:wrapcheck
			}
			msg.Result = result
		} else {
			// It's a notification
			return t.conn.Notify(ctx, msg.Method, msg.Params) //nolint:wrapcheck
		}
		return nil
	}

	// Manual implementation for raw reader/writer
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	// Write Content-Length header
	header := fmt.Sprintf("Content-Length: %d\r\n\r\n", len(data))
	if _, err := t.writer.Write([]byte(header)); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	// Write message body
	if _, err := t.writer.Write(data); err != nil {
		return fmt.Errorf("failed to write body: %w", err)
	}

	if t.logger != nil {
		t.logger.Debug("Sent message", "method", msg.Method, "id", msg.ID)
	}

	return nil
}

// Receive 从语言服务器接收消息。
func (t *Transport) Receive(_ context.Context) (*Message, error) {
	if t.conn != nil {
		// This is handled by the connection's handler
		return nil, fmt.Errorf("receive not supported with existing connection")
	}

	// Read headers
	headers := make(map[string]string)
	scanner := bufio.NewScanner(t.reader)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}

		parts := strings.SplitN(line, ": ", 2)
		if len(parts) == 2 {
			headers[parts[0]] = parts[1]
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read headers: %w", err)
	}

	// Get content length
	contentLengthStr, ok := headers["Content-Length"]
	if !ok {
		return nil, fmt.Errorf("missing Content-Length header")
	}

	contentLength, err := strconv.Atoi(contentLengthStr)
	if err != nil {
		return nil, fmt.Errorf("invalid Content-Length: %w", err)
	}

	// Read body
	body := make([]byte, contentLength)
	if _, err := io.ReadFull(t.reader, body); err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}

	// Parse message
	var msg Message
	if err := json.Unmarshal(body, &msg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal message: %w", err)
	}

	if t.logger != nil {
		t.logger.Debug("Received message", "method", msg.Method, "id", msg.ID)
	}

	return &msg, nil
}

// Close 关闭传输。
func (t *Transport) Close() error {
	if t.conn != nil {
		// Connection will be closed by the client
		return nil
	}

	if closer, ok := t.writer.(io.Closer); ok {
		if err := closer.Close(); err != nil {
			return err //nolint:wrapcheck
		}
	}

	if closer, ok := t.reader.(io.Closer); ok {
		if err := closer.Close(); err != nil {
			return err //nolint:wrapcheck
		}
	}

	return nil
}

// StreamTransport 为 JSON-RPC 通信提供双向流。
type StreamTransport struct {
	reader io.Reader
	writer io.Writer
	closer io.Closer
}

// NewStreamTransport 创建一个新的流传输。
func NewStreamTransport(reader io.Reader, writer io.Writer, closer io.Closer) *StreamTransport {
	return &StreamTransport{
		reader: reader,
		writer: writer,
		closer: closer,
	}
}

// Read 实现 io.Reader 接口。
func (s *StreamTransport) Read(p []byte) (n int, err error) {
	return s.reader.Read(p) //nolint:wrapcheck
}

// Write 实现 io.Writer 接口。
func (s *StreamTransport) Write(p []byte) (n int, err error) {
	return s.writer.Write(p) //nolint:wrapcheck
}

// Close 实现 io.Closer 接口。
func (s *StreamTransport) Close() error {
	if s.closer != nil {
		return s.closer.Close() //nolint:wrapcheck
	}
	return nil
}

// ObjectStream 从传输创建一个 jsonrpc2.ObjectStream。
func (s *StreamTransport) ObjectStream() jsonrpc2.ObjectStream {
	return jsonrpc2.NewBufferedStream(s, jsonrpc2.VSCodeObjectCodec{})
}
