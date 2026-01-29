//go:build windows

package desktop

import (
	"syscall"
	"unsafe"
)

var (
	kernel32              = syscall.NewLazyDLL("kernel32.dll")
	shell32               = syscall.NewLazyDLL("shell32.dll")
	user32Icon            = syscall.NewLazyDLL("user32.dll")
	procSendMessage       = user32Icon.NewProc("SendMessageW")
	procGetModuleFileName = kernel32.NewProc("GetModuleFileNameW")
	procExtractIcon       = shell32.NewProc("ExtractIconW")
)

const (
	WM_SETICON = 0x0080
	ICON_SMALL = 0
	ICON_BIG   = 1
)

func SetWindowIcon(hwnd unsafe.Pointer) {
	// robust strategy: Extract the icon from the executable file itself.
	// This grabs specific index 0 (the main icon).

	// 1. Get Path to Exe
	buf := make([]uint16, 260)
	procGetModuleFileName.Call(0, uintptr(unsafe.Pointer(&buf[0])), 260)

	// 2. Extract Icon
	iconHandle, _, _ := procExtractIcon.Call(
		0,
		uintptr(unsafe.Pointer(&buf[0])),
		0, // Index 0
	)

	if iconHandle == 0 || iconHandle == 1 {
		return
	}

	// 3. Set Icon
	procSendMessage.Call(uintptr(hwnd), WM_SETICON, ICON_BIG, iconHandle)
	procSendMessage.Call(uintptr(hwnd), WM_SETICON, ICON_SMALL, iconHandle)
}
