package pony

import (
	"fmt"
	"strconv"
	"strings"
)

// SizeConstraint 表示带有单位的尺寸约束。
type SizeConstraint struct {
	value int
	unit  string // "", "%", "auto", "min", "max"
}

// parseSizeConstraint 解析尺寸字符串，如 "50%", "20", "auto"。
func parseSizeConstraint(s string) SizeConstraint {
	s = strings.TrimSpace(s)
	if s == "" {
		return SizeConstraint{unit: UnitAuto}
	}

	// 检查特殊关键字
	switch s {
	case UnitAuto:
		return SizeConstraint{unit: UnitAuto}
	case UnitMin:
		return SizeConstraint{unit: UnitMin}
	case UnitMax:
		return SizeConstraint{unit: UnitMax}
	}

	// 检查百分比
	if strings.HasSuffix(s, UnitPercent) {
		valStr := strings.TrimSuffix(s, UnitPercent)
		if val, err := strconv.Atoi(valStr); err == nil {
			return SizeConstraint{value: val, unit: UnitPercent}
		}
	}

	// 固定尺寸
	if val, err := strconv.Atoi(s); err == nil {
		return SizeConstraint{value: val, unit: ""}
	}

	// 无效，默认为自动
	return SizeConstraint{unit: UnitAuto}
}

// Apply 应用尺寸约束以获取实际尺寸。
func (sc SizeConstraint) Apply(available, content int) int {
	// 零值（unit="" 且 value=0）被视为自动
	if sc.unit == "" && sc.value == 0 {
		if content > available {
			return available
		}
		return content
	}

	switch sc.unit {
	case "%":
		// 可用空间的百分比
		result := available * sc.value / 100
		if result < 0 {
			return 0
		}
		if result > available {
			return available
		}
		return result

	case UnitAuto:
		// 基于内容的尺寸
		if content > available {
			return available
		}
		return content

	case UnitMin:
		// 最小尺寸（内容或 0）
		if content < 0 {
			return 0
		}
		return content

	case UnitMax:
		// 最大可用尺寸
		return available

	default:
		// 固定尺寸
		if sc.value < 0 {
			return 0
		}
		if sc.value > available {
			return available
		}
		return sc.value
	}
}

// String 返回字符串表示形式。
func (sc SizeConstraint) String() string {
	switch sc.unit {
	case UnitPercent:
		return fmt.Sprintf("%d%%", sc.value)
	case UnitAuto, UnitMin, UnitMax:
		return sc.unit
	default:
		return fmt.Sprintf("%d", sc.value)
	}
}

// IsAuto 如果这是自动约束，则返回 true。
func (sc SizeConstraint) IsAuto() bool {
	// 零值（unit="" 且 value=0）是自动
	// 否则 unit 必须明确为 "auto"
	return (sc.unit == "" && sc.value == 0) || sc.unit == UnitAuto
}

// IsFixed 如果这是固定尺寸约束，则返回 true。
func (sc SizeConstraint) IsFixed() bool {
	return sc.unit == ""
}

// IsPercent 如果这是百分比约束，则返回 true。
func (sc SizeConstraint) IsPercent() bool {
	return sc.unit == UnitPercent
}

// NewFixedConstraint 创建固定尺寸约束。
func NewFixedConstraint(size int) SizeConstraint {
	return SizeConstraint{value: size, unit: ""}
}

// NewPercentConstraint 创建百分比约束。
func NewPercentConstraint(percent int) SizeConstraint {
	return SizeConstraint{value: percent, unit: UnitPercent}
}
