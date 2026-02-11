package protocol

import (
	"fmt"
	"log/slog"
)

// PatternInfo 是表示 glob 模式的类型的接口。
type PatternInfo interface {
	GetPattern() string
	GetBasePath() string
	isPattern() // marker method
}

// StringPattern 为字符串模式实现了 PatternInfo 接口。
type StringPattern struct {
	Pattern string
}

// GetPattern 返回 glob 模式字符串。
func (p StringPattern) GetPattern() string { return p.Pattern }

// GetBasePath 为简单模式返回空字符串。
func (p StringPattern) GetBasePath() string { return "" }
func (p StringPattern) isPattern()          {}

// RelativePatternInfo 为 RelativePattern 实现了 PatternInfo 接口。
type RelativePatternInfo struct {
	RP       RelativePattern
	BasePath string
}

// GetPattern 返回 glob 模式字符串。
func (p RelativePatternInfo) GetPattern() string { return p.RP.Pattern }

// GetBasePath 返回模式的基础路径。
func (p RelativePatternInfo) GetBasePath() string { return p.BasePath }
func (p RelativePatternInfo) isPattern()          {}

// AsPattern 将 GlobPattern 转换为 PatternInfo 对象。
func (g *GlobPattern) AsPattern() (PatternInfo, error) {
	if g.Value == nil {
		return nil, fmt.Errorf("nil pattern")
	}

	var err error

	switch v := g.Value.(type) {
	case string:
		return StringPattern{Pattern: v}, nil

	case RelativePattern:
		// Handle BaseURI which could be string or DocumentUri
		basePath := ""
		switch baseURI := v.BaseURI.Value.(type) {
		case string:
			basePath, err = DocumentURI(baseURI).Path()
			if err != nil {
				slog.Error("Failed to convert URI to path", "uri", baseURI, "error", err)
				return nil, fmt.Errorf("invalid URI: %s", baseURI)
			}

		case DocumentURI:
			basePath, err = baseURI.Path()
			if err != nil {
				slog.Error("Failed to convert DocumentURI to path", "uri", baseURI, "error", err)
				return nil, fmt.Errorf("invalid DocumentURI: %s", baseURI)
			}

		default:
			return nil, fmt.Errorf("unknown BaseURI type: %T", v.BaseURI.Value)
		}

		return RelativePatternInfo{RP: v, BasePath: basePath}, nil

	default:
		return nil, fmt.Errorf("unknown pattern type: %T", g.Value)
	}
}
