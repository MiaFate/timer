//go:build windows

package desktop

import (
	"syscall"
	"timer/internal/domain"
	"unsafe"

	"golang.org/x/sys/windows"
)

type WindowsService struct{}

func NewWindowService() domain.WindowService {
	return &WindowsService{}
}

var (
	user32                       = windows.NewLazySystemDLL("user32.dll")
	procGetWindowTextW           = user32.NewProc("GetWindowTextW")
	procGetWindowTextLengthW     = user32.NewProc("GetWindowTextLengthW")
	procIsWindowVisible          = user32.NewProc("IsWindowVisible")
	procGetWindowThreadProcessId = user32.NewProc("GetWindowThreadProcessId")
	procGetForegroundWindow      = user32.NewProc("GetForegroundWindow")
)

func (s *WindowsService) GetOpenWindows() ([]domain.AppInfo, error) {
	var apps []domain.AppInfo
	seenPIDs := make(map[uint32]bool)

	cb := syscall.NewCallback(func(hwnd windows.Handle, lparam uintptr) uintptr {
		if isVisible, _, _ := procIsWindowVisible.Call(uintptr(hwnd)); isVisible == 0 {
			return 1 // Continue
		}

		pid := getPid(hwnd)

		if seenPIDs[pid] {
			return 1
		}

		// Filter out empty titles
		title := getWindowText(hwnd)
		if title == "" {
			return 1
		}

		seenPIDs[pid] = true
		apps = append(apps, domain.AppInfo{
			PID:   pid,
			Title: title,
			HWnd:  hwnd,
		})
		return 1 // Continue
	})

	windows.EnumWindows(cb, nil)
	return apps, nil
}

func (s *WindowsService) GetForegroundWindowPID() (uint32, error) {
	hwnd, _, _ := procGetForegroundWindow.Call()
	if hwnd == 0 {
		return 0, nil
	}

	pid := getPid(windows.Handle(hwnd))
	return pid, nil
}

func getPid(hwnd windows.Handle) uint32 {
	var pid uint32
	procGetWindowThreadProcessId.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&pid)))
	return pid
}

func getWindowText(hwnd windows.Handle) string {
	len, _, _ := procGetWindowTextLengthW.Call(uintptr(hwnd))
	if len == 0 {
		return ""
	}
	buf := make([]uint16, len+1)
	procGetWindowTextW.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&buf[0])), uintptr(len+1))
	return windows.UTF16ToString(buf)
}
