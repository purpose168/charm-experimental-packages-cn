// Package protocol 为语言服务器协议（Language Server Protocol，LSP）提供类型和函数。
package protocol

import "fmt"

// WorkspaceSymbolResult 是表示工作区符号的类型的接口。
type WorkspaceSymbolResult interface {
	GetName() string
	GetLocation() Location
	isWorkspaceSymbol() // marker method
}

// GetName 返回符号名称。
func (ws *WorkspaceSymbol) GetName() string { return ws.Name }

// GetLocation 返回符号位置。
func (ws *WorkspaceSymbol) GetLocation() Location {
	switch v := ws.Location.Value.(type) {
	case Location:
		return v
	case LocationUriOnly:
		return Location{URI: v.URI}
	}
	return Location{}
}
func (ws *WorkspaceSymbol) isWorkspaceSymbol() {}

// GetName returns the symbol name.
func (si *SymbolInformation) GetName() string { return si.Name }

// GetLocation returns the symbol location.
func (si *SymbolInformation) GetLocation() Location { return si.Location }
func (si *SymbolInformation) isWorkspaceSymbol()    {}

// Results 将 Value 转换为 WorkspaceSymbolResult 切片。
func (r Or_Result_workspace_symbol) Results() ([]WorkspaceSymbolResult, error) {
	if r.Value == nil {
		return make([]WorkspaceSymbolResult, 0), nil
	}
	switch v := r.Value.(type) {
	case []WorkspaceSymbol:
		results := make([]WorkspaceSymbolResult, len(v))
		for i := range v {
			results[i] = &v[i]
		}
		return results, nil
	case []SymbolInformation:
		results := make([]WorkspaceSymbolResult, len(v))
		for i := range v {
			results[i] = &v[i]
		}
		return results, nil
	default:
		return nil, fmt.Errorf("unknown symbol type: %T", r.Value)
	}
}

// DocumentSymbolResult 是表示文档符号的类型的接口。
type DocumentSymbolResult interface {
	GetRange() Range
	GetName() string
	isDocumentSymbol() // marker method
}

// GetRange 返回符号范围。
func (ds *DocumentSymbol) GetRange() Range { return ds.Range }

// GetName 返回符号名称。
func (ds *DocumentSymbol) GetName() string   { return ds.Name }
func (ds *DocumentSymbol) isDocumentSymbol() {}

// GetRange 从其位置返回符号范围。
func (si *SymbolInformation) GetRange() Range { return si.Location.Range }

// Note: SymbolInformation already has GetName() implemented above.
func (si *SymbolInformation) isDocumentSymbol() {}

// Results 将 Value 转换为 DocumentSymbolResult 切片。
func (r Or_Result_textDocument_documentSymbol) Results() ([]DocumentSymbolResult, error) {
	if r.Value == nil {
		return make([]DocumentSymbolResult, 0), nil
	}
	switch v := r.Value.(type) {
	case []DocumentSymbol:
		results := make([]DocumentSymbolResult, len(v))
		for i := range v {
			results[i] = &v[i]
		}
		return results, nil
	case []SymbolInformation:
		results := make([]DocumentSymbolResult, len(v))
		for i := range v {
			results[i] = &v[i]
		}
		return results, nil
	default:
		return nil, fmt.Errorf("unknown document symbol type: %T", v)
	}
}

// TextEditResult 是可以用作文本编辑的类型的接口。
type TextEditResult interface {
	GetRange() Range
	GetNewText() string
	isTextEdit() // marker method
}

// GetRange 返回编辑范围。
func (te *TextEdit) GetRange() Range { return te.Range }

// GetNewText 返回编辑的新文本。
func (te *TextEdit) GetNewText() string { return te.NewText }
func (te *TextEdit) isTextEdit()        {}

// AsTextEdit 将 Or_TextDocumentEdit_edits_Elem 转换为 TextEdit。
func (e Or_TextDocumentEdit_edits_Elem) AsTextEdit() (TextEdit, error) {
	if e.Value == nil {
		return TextEdit{}, fmt.Errorf("nil text edit")
	}
	switch v := e.Value.(type) {
	case TextEdit:
		return v, nil
	case AnnotatedTextEdit:
		return TextEdit{
			Range:   v.Range,
			NewText: v.NewText,
		}, nil
	default:
		return TextEdit{}, fmt.Errorf("unknown text edit type: %T", e.Value)
	}
}
