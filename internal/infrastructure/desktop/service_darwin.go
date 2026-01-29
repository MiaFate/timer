//go:build darwin

package desktop

import (
	"os/exec"
	"strconv"
	"strings"
	"timer/internal/domain"
)

type DarwinService struct{}

func NewWindowService() domain.WindowService {
	return &DarwinService{}
}

func (s *DarwinService) GetOpenWindows() ([]domain.AppInfo, error) {
	// macOS doesn't have a simple way to list ALL open windows with PIDs like EnumWindows without Cgo/CoreGraphics
	// For now, we return empty or just the frontmost to satisfy the interface.
	// Listing all app names:
	out, err := exec.Command("osascript", "-e", `tell application "System Events" to get name of every application process whose background only is false`).Output()
	if err != nil {
		return []domain.AppInfo{}, nil
	}

	names := strings.Split(string(out), ", ")
	var apps []domain.AppInfo
	for i, name := range names {
		apps = append(apps, domain.AppInfo{
			PID:   uint32(i + 1000), // Dummy PID for listing if we can't get it easily
			Title: strings.TrimSpace(name),
		})
	}

	return apps, nil
}

func (s *DarwinService) GetForegroundWindowPID() (uint32, error) {
	// Strategy: Use AppleScript to get the unix id of the frontmost process
	out, err := exec.Command("osascript", "-e", `tell application "System Events" to get unix id of first application process whose frontmost is true`).Output()
	if err != nil {
		return 0, nil
	}

	pidStr := strings.TrimSpace(string(out))
	pid, err := strconv.ParseUint(pidStr, 10, 32)
	if err != nil {
		return 0, nil
	}

	return uint32(pid), nil
}
