package ansi

// Modes 表示可以设置或重置的终端模式。默认情况下，
// 所有模式都是 [ModeNotRecognized]。
type Modes map[Mode]ModeSetting

// Get 返回终端模式的设置。如果模式未设置，则返回 [ModeNotRecognized]。
func (m Modes) Get(mode Mode) ModeSetting {
	return m[mode]
}

// Delete 删除终端模式。这与将模式设置为 [ModeNotRecognized] 的效果相同。
func (m Modes) Delete(mode Mode) {
	delete(m, mode)
}

// Set 将终端模式设置为 [ModeSet]。
func (m Modes) Set(modes ...Mode) {
	for _, mode := range modes {
		m[mode] = ModeSet
	}
}

// PermanentlySet 将终端模式设置为 [ModePermanentlySet]。
func (m Modes) PermanentlySet(modes ...Mode) {
	for _, mode := range modes {
		m[mode] = ModePermanentlySet
	}
}

// Reset 将终端模式设置为 [ModeReset]。
func (m Modes) Reset(modes ...Mode) {
	for _, mode := range modes {
		m[mode] = ModeReset
	}
}

// PermanentlyReset 将终端模式设置为 [ModePermanentlyReset]。
func (m Modes) PermanentlyReset(modes ...Mode) {
	for _, mode := range modes {
		m[mode] = ModePermanentlyReset
	}
}

// IsSet 如果模式设置为 [ModeSet] 或 [ModePermanentlySet]，则返回 true。
func (m Modes) IsSet(mode Mode) bool {
	return m[mode].IsSet()
}

// IsPermanentlySet 如果模式设置为 [ModePermanentlySet]，则返回 true。
func (m Modes) IsPermanentlySet(mode Mode) bool {
	return m[mode].IsPermanentlySet()
}

// IsReset 如果模式设置为 [ModeReset] 或 [ModePermanentlyReset]，则返回 true。
func (m Modes) IsReset(mode Mode) bool {
	return m[mode].IsReset()
}

// IsPermanentlyReset 如果模式设置为 [ModePermanentlyReset]，则返回 true。
func (m Modes) IsPermanentlyReset(mode Mode) bool {
	return m[mode].IsPermanentlyReset()
}
