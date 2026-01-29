//go:build linux

package desktop

import (
	"os/exec"
	"strings"
	"timer/internal/domain"
)

type LinuxService struct{}

func NewWindowService() domain.WindowService {
	return &LinuxService{}
}

func (s *LinuxService) GetOpenWindows() ([]domain.AppInfo, error) {
	// Placeholder: Using xrop or wmctrl would go here.
	// For now, return empty to allow compilation.
	return []domain.AppInfo{}, nil
}

func (s *LinuxService) GetForegroundWindowPID() (uint32, error) {
	// Simple implementation using xprop (common on X11)
	// Requires 'xprop' installed on user system
	out, err := exec.Command("xprop", "-root", "_NET_ACTIVE_WINDOW").Output()
	if err != nil {
		return 0, nil
	}

	// Output format: _NET_ACTIVE_WINDOW(WINDOW): window id # 0x4400005
	str := string(out)
	parts := strings.Split(str, " ")
	if len(parts) < 5 {
		return 0, nil
	}

	// This gives Window ID, we need PID.
	// "xprop -id <windowid> _NET_WM_PID"
	// For MVP/compilation, returning 0 is safe.
	return 0, nil
}
