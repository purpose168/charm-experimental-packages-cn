package ansi

import (
	"sync"

	"github.com/purpose168/charm-experimental-packages-cn/ansi/parser"
)

var parserPool = sync.Pool{
	New: func() any {
		p := NewParser()
		p.SetParamsSize(parser.MaxParamsSize)
		p.SetDataSize(1024 * 4) // 4KB of data buffer
		return p
	},
}

// GetParser 从同步池中获取一个解析器。
func GetParser() *Parser {
	return parserPool.Get().(*Parser)
}

// PutParser 将解析器返回给同步池。解析器会自动重置。
func PutParser(p *Parser) {
	p.Reset()
	p.dataLen = 0
	parserPool.Put(p)
}
