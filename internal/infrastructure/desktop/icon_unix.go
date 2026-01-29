//go:build linux || darwin

package desktop

import "unsafe"

func SetWindowIcon(hwnd unsafe.Pointer) {
	// No-op for now on Linux/Mac
	// Linux uses .desktop files or X11 properties
	// Mac uses bundle resources
}
