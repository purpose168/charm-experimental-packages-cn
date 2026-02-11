//go:build !windows
// +build !windows

package main

import (
	"os"
	"os/signal"
	"syscall"
)

// listenForResize 监听终端窗口大小变化信号
// 当接收到 SIGWINCH 信号时，调用传入的回调函数
func listenForResize(fn func()) {
	// 创建一个信号通道，用于接收系统信号
	sig := make(chan os.Signal, 1)
	// 注册监听 SIGWINCH 信号（窗口大小变化信号）
	signal.Notify(sig, syscall.SIGWINCH)

	// 持续监听信号通道，当接收到信号时调用回调函数
	for range sig {
		fn()
	}
}
