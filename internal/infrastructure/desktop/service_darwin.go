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
	// Strategy: Use AppleScript to get PID and name together
	// Format: pid1, name1, pid2, name2...
	script := `set output to ""
	tell application "System Events"
		set procs to every application process whose background only is false
		repeat with p in procs
			set output to output & (unix id of p as text) & "|" & (name of p as text) & "\n"
		end repeat
	end tell
	return output`

	out, err := exec.Command("osascript", "-e", script).Output()
	if err != nil {
		return []domain.AppInfo{}, nil
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	var apps []domain.AppInfo
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.Split(line, "|")
		if len(parts) < 2 {
			continue
		}

		pid64, _ := strconv.ParseUint(strings.TrimSpace(parts[0]), 10, 32)
		apps = append(apps, domain.AppInfo{
			PID:   uint32(pid64),
			Title: strings.TrimSpace(parts[1]),
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
