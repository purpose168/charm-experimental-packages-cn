// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package protocol

import (
	"encoding/json"
	"fmt"
)

// DocumentChange 是各种文件编辑操作的联合类型。
//
// 此结构体中恰好有一个字段非空；请参阅 [DocumentChange.Valid]。
//
// 参见 https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#resourceChanges
type DocumentChange struct {
	TextDocumentEdit *TextDocumentEdit
	CreateFile       *CreateFile
	RenameFile       *RenameFile
	DeleteFile       *DeleteFile
}

// Valid 报告 DocumentChange 联合类型值是否有效，
// 即恰好有一个创建、删除、编辑或重命名操作。
func (d DocumentChange) Valid() bool {
	n := 0
	if d.TextDocumentEdit != nil {
		n++
	}
	if d.CreateFile != nil {
		n++
	}
	if d.RenameFile != nil {
		n++
	}
	if d.DeleteFile != nil {
		n++
	}
	return n == 1
}

// UnmarshalJSON 实现 json.Unmarshaler 接口。
func (d *DocumentChange) UnmarshalJSON(data []byte) error {
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return err //nolint:wrapcheck
	}

	if _, ok := m["textDocument"]; ok {
		d.TextDocumentEdit = new(TextDocumentEdit)
		return json.Unmarshal(data, d.TextDocumentEdit) //nolint:wrapcheck
	}

	// The {Create,Rename,Delete}File types all share a 'kind' field.
	kind := m["kind"]
	switch kind {
	case "create":
		d.CreateFile = new(CreateFile)
		return json.Unmarshal(data, d.CreateFile) //nolint:wrapcheck
	case "rename":
		d.RenameFile = new(RenameFile)
		return json.Unmarshal(data, d.RenameFile) //nolint:wrapcheck
	case "delete":
		d.DeleteFile = new(DeleteFile)
		return json.Unmarshal(data, d.DeleteFile) //nolint:wrapcheck
	}
	return fmt.Errorf("DocumentChanges: unexpected kind: %q", kind)
}

// MarshalJSON 实现 json.Marshaler 接口。
func (d *DocumentChange) MarshalJSON() ([]byte, error) {
	if d.TextDocumentEdit != nil {
		return json.Marshal(d.TextDocumentEdit) //nolint:wrapcheck
	} else if d.CreateFile != nil {
		return json.Marshal(d.CreateFile) //nolint:wrapcheck
	} else if d.RenameFile != nil {
		return json.Marshal(d.RenameFile) //nolint:wrapcheck
	} else if d.DeleteFile != nil {
		return json.Marshal(d.DeleteFile) //nolint:wrapcheck
	}
	return nil, fmt.Errorf("empty DocumentChanges union value")
}
