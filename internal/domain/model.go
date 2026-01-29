package domain

import (
	"golang.org/x/sys/windows"
)

type AppInfo struct {
	PID   uint32         `json:"pid"`
	Title string         `json:"title"`
	HWnd  windows.Handle `json:"hwnd"`
}

type WindowService interface {
	GetOpenWindows() ([]AppInfo, error)
	GetForegroundWindowPID() (uint32, error)
}
