//go:build windows
// +build windows

package main

// listenForResize 监听终端窗口大小变化
// 在 Windows 平台上，此函数为空实现
func listenForResize(func()) {}
