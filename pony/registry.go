package pony

import (
	"sync"
)

// ComponentFactory 是一个从属性和子元素创建 Element 的函数。
// 自定义组件应该实现此签名。
type ComponentFactory func(props Props, children []Element) Element

var (
	registry   = make(map[string]ComponentFactory)
	registryMu sync.RWMutex
)

// Register 使用给定名称注册自定义组件。
// 工厂函数将被调用来创建组件实例。
//
// 示例：
//
//	pony.Register("badge", func(props Props, children []Element) Element {
//	    return &Badge{
//	        text:  props.Get("text"),
//	        color: props.Get("color"),
//	    }
//	})
func Register(name string, factory ComponentFactory) {
	registryMu.Lock()
	defer registryMu.Unlock()
	registry[name] = factory
}

// Unregister 移除已注册的组件。
func Unregister(name string) {
	registryMu.Lock()
	defer registryMu.Unlock()
	delete(registry, name)
}

// GetComponent 通过名称检索组件工厂。
// 如果组件未注册，返回 nil。
func GetComponent(name string) (ComponentFactory, bool) {
	registryMu.RLock()
	defer registryMu.RUnlock()
	factory, ok := registry[name]
	return factory, ok
}

// RegisteredComponents 返回所有已注册组件名称的列表。
func RegisteredComponents() []string {
	registryMu.RLock()
	defer registryMu.RUnlock()

	names := make([]string, 0, len(registry))
	for name := range registry {
		names = append(names, name)
	}
	return names
}

// ClearRegistry 移除所有已注册的组件。
// 对测试很有用。
func ClearRegistry() {
	registryMu.Lock()
	defer registryMu.Unlock()
	registry = make(map[string]ComponentFactory)
}
