package lsp

import (
	"cmp"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/purpose168/charm-experimental-packages-cn/powernap/pkg/lsp/protocol"
)

// TextDocumentSyncManager 管理与语言服务器的文本文档同步。
type TextDocumentSyncManager struct {
	client    *Client
	documents map[string]*Document
	syncKind  protocol.TextDocumentSyncKind
}

// Document 表示一个打开的文本文档。
type Document struct {
	URI        string
	LanguageID string
	Version    int
	Content    string
}

// NewTextDocumentSyncManager 创建一个新的文本文档同步管理器。
func NewTextDocumentSyncManager(client *Client) *TextDocumentSyncManager {
	syncKind := protocol.Full // Default to full sync

	// Extract sync kind from capabilities
	switch v := client.capabilities.TextDocumentSync.(type) {
	case float64:
		syncKind = protocol.TextDocumentSyncKind(int(v)) //nolint:gosec
	case int:
		syncKind = protocol.TextDocumentSyncKind(v) //nolint:gosec
	case map[string]any:
		// It's a TextDocumentSyncOptions object
		if change, ok := v["change"].(float64); ok {
			syncKind = protocol.TextDocumentSyncKind(int(change)) //nolint:gosec
		}
	case *protocol.TextDocumentSyncOptions:
		syncKind = v.Change
	}

	return &TextDocumentSyncManager{
		client:    client,
		documents: make(map[string]*Document),
		syncKind:  syncKind,
	}
}

// Open 打开一个新的文本文档。
func (m *TextDocumentSyncManager) Open(uri, languageID, content string) error {
	if _, exists := m.documents[uri]; exists {
		return fmt.Errorf("document already open: %s", uri)
	}

	doc := &Document{
		URI:        uri,
		LanguageID: languageID,
		Version:    1,
		Content:    content,
	}

	m.documents[uri] = doc

	return m.client.NotifyDidOpenTextDocument(m.client.ctx, uri, languageID, doc.Version, content)
}

// Change 对打开的文档应用更改。
func (m *TextDocumentSyncManager) Change(uri string, changes []protocol.TextDocumentContentChangeEvent) error {
	doc, exists := m.documents[uri]
	if !exists {
		return fmt.Errorf("document not open: %s", uri)
	}

	// Apply changes based on sync kind
	switch m.syncKind {
	case protocol.None:
		// Server doesn't want document change notifications
		return nil
	case protocol.Full:
		// For full sync, we expect a single change with the full content
		if len(changes) > 0 {
			m.applyFullDocumentChange(doc, changes[0])
		}
	case protocol.Incremental:
		// Apply incremental changes
		for _, change := range changes {
			if err := m.applyIncrementalChange(doc, change); err != nil {
				return err
			}
		}
	}

	doc.Version++

	return m.client.NotifyDidChangeTextDocument(m.client.ctx, uri, doc.Version, changes)
}

// Close 关闭一个文本文档。
func (m *TextDocumentSyncManager) Close(uri string) error {
	if _, exists := m.documents[uri]; !exists {
		return fmt.Errorf("document not open: %s", uri)
	}

	delete(m.documents, uri)

	// Send didClose notification
	params := map[string]any{
		"textDocument": map[string]any{
			"uri": uri,
		},
	}

	return m.client.conn.Notify(m.client.ctx, MethodTextDocumentDidClose, params) //nolint:wrapcheck
}

// Save 通知服务器文档已保存。
func (m *TextDocumentSyncManager) Save(uri string, includeText bool) error {
	doc, exists := m.documents[uri]
	if !exists {
		return fmt.Errorf("document not open: %s", uri)
	}

	params := map[string]any{
		"textDocument": map[string]any{
			"uri": uri,
		},
	}

	if includeText {
		params["text"] = doc.Content
	}

	return m.client.conn.Notify(m.client.ctx, MethodTextDocumentDidSave, params) //nolint:wrapcheck
}

// GetDocument 返回给定 URI 的文档。
func (m *TextDocumentSyncManager) GetDocument(uri string) (*Document, bool) {
	doc, exists := m.documents[uri]
	return doc, exists
}

func (m *TextDocumentSyncManager) applyFullDocumentChange(doc *Document, change protocol.TextDocumentContentChangeEvent) {
	full, ok := change.Value.(protocol.TextDocumentContentChangeWholeDocument)
	if !ok {
		return
	}
	doc.Content = full.Text
}

// applyIncrementalChange 对文档应用增量更改。
func (m *TextDocumentSyncManager) applyIncrementalChange(doc *Document, change protocol.TextDocumentContentChangeEvent) error {
	partial, ok := change.Value.(protocol.TextDocumentContentChangePartial)
	if !ok {
		return nil
	}
	if partial.Range == nil {
		// Full document change
		doc.Content = partial.Text
		return nil
	}

	// Convert content to lines for easier manipulation
	lines := strings.Split(doc.Content, "\n")
	lineC := uint32(len(lines)) //nolint:gosec

	// Validate range
	if partial.Range.Start.Line >= lineC {
		return fmt.Errorf("invalid start line: %d", partial.Range.Start.Line)
	}
	if partial.Range.End.Line >= lineC {
		return fmt.Errorf("invalid end line: %d", partial.Range.End.Line)
	}

	// Calculate the start and end positions in the document
	var startPos uint32
	for i := uint32(0); i < partial.Range.Start.Line; i++ {
		startPos += uint32(len(lines[i])) + 1 //nolint:gosec
	}
	startPos += partial.Range.Start.Character

	var endPos uint32
	for i := uint32(0); i < partial.Range.End.Line; i++ {
		endPos += uint32(len(lines[i])) + 1 //nolint:gosec
	}
	endPos += partial.Range.End.Character

	// Apply the partial
	newContent := doc.Content[:startPos] + partial.Text + doc.Content[endPos:]
	doc.Content = newContent

	return nil
}

// CreateFullDocumentChange 为完整文档同步创建更改事件。
func CreateFullDocumentChange(content string) []protocol.TextDocumentContentChangeEvent {
	return []protocol.TextDocumentContentChangeEvent{
		{
			Value: protocol.TextDocumentContentChangeWholeDocument{
				Text: content,
			},
		},
	}
}

// OpenFile 从文件路径打开一个新的文本文档。
func (m *TextDocumentSyncManager) OpenFile(filePath, content string) error {
	// Convert to absolute path
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Detect language from file path
	languageID := cmp.Or(string(DetectLanguage(absPath)), "plaintext")
	// Create file URI
	uri := FilePathToURI(absPath)

	return m.Open(uri, languageID, content)
}

// FilePathToURI 将文件路径转换为文件 URI。
//
// 已弃用：使用 [protocol.URIFromPath]。
func FilePathToURI(path string) string {
	return string(protocol.URIFromPath(path))
}

// URIToFilePath 将文件 URI 转换为文件路径。
func URIToFilePath(uri string) string {
	// Remove file:// prefix
	path := strings.TrimPrefix(uri, "file://")
	path = strings.TrimPrefix(path, "file:///")

	// Convert to native path separators
	return filepath.FromSlash(path)
}
