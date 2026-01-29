//go:build linux

package desktop

import (
	"os/exec"
	"strconv"
	"strings"
	"timer/internal/domain"
)

type LinuxService struct{}

func NewWindowService() domain.WindowService {
	return &LinuxService{}
}

func (s *LinuxService) GetOpenWindows() ([]domain.AppInfo, error) {
	// Use wmctrl to list windows if available
	out, err := exec.Command("wmctrl", "-lp").Output()
	if err != nil {
		// Fallback: return empty
		return []domain.AppInfo{}, nil
	}

	lines := strings.Split(string(out), "\n")
	var apps []domain.AppInfo
	seenPIDs := make(map[uint32]bool)

	for _, line := range lines {
		parts := strings.Fields(line)
		// Format: 0x03a00003  0 2564   hostname Title
		if len(parts) < 5 {
			continue
		}

		pid64, _ := strconv.ParseUint(parts[2], 10, 32)
		pid := uint32(pid64)
		if pid == 0 || seenPIDs[pid] {
			continue
		}

		title := strings.Join(parts[4:], " ")
		seenPIDs[pid] = true
		apps = append(apps, domain.AppInfo{
			PID:   pid,
			Title: title,
		})
	}
	return apps, nil
}

func (s *LinuxService) GetForegroundWindowPID() (uint32, error) {
	// Strategy: Get active window ID, then get its PID
	out, err := exec.Command("xprop", "-root", "_NET_ACTIVE_WINDOW").Output()
	if err != nil {
		return 0, nil
	}

	// Output: _NET_ACTIVE_WINDOW(WINDOW): window id # 0x4400005
	parts := strings.Split(strings.TrimSpace(string(out)), " ")
	windowID := parts[len(parts)-1]

	if windowID == "0x0" {
		return 0, nil
	}

	pidOut, err := exec.Command("xprop", "-id", windowID, "_NET_WM_PID").Output()
	if err != nil {
		return 0, nil
	}

	// Output: _NET_WM_PID(CARDINAL) = 2564
	pidParts := strings.Split(strings.TrimSpace(string(pidOut)), " ")
	pidStr := pidParts[len(pidParts)-1]

	pid, _ := strconv.ParseUint(pidStr, 10, 32)
	return uint32(pid), nil
}
