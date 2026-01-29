package domain

type AppInfo struct {
	PID   uint32 `json:"pid"`
	Title string `json:"title"`
}

type WindowService interface {
	GetOpenWindows() ([]AppInfo, error)
	GetForegroundWindowPID() (uint32, error)
}
