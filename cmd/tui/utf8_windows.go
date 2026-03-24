//go:build windows

package main

import "syscall"

func setUTF8Mode() {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	setConsoleOutputCP := kernel32.NewProc("SetConsoleOutputCP")
	setConsoleCP := kernel32.NewProc("SetConsoleCP")

	setConsoleOutputCP.Call(uintptr(65001))
	setConsoleCP.Call(uintptr(65001))
}
