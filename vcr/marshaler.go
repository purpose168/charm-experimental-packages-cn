package vcr

import (
	"bytes"
	"fmt"

	"go.yaml.in/yaml/v4"
)

// customMarshaler 将数据编码为YAML格式，使用2个空格的缩进和紧凑的序列缩进。
func customMarshaler(in any) ([]byte, error) {
	var buff bytes.Buffer
	enc := yaml.NewEncoder(&buff)
	enc.SetIndent(2) // 设置缩进为2个空格
	enc.CompactSeqIndent() // 使用紧凑的序列缩进
	if err := enc.Encode(in); err != nil {
		return nil, fmt.Errorf("vcr: unable to encode to yaml: %w", err)
	}
	return buff.Bytes(), nil
}
