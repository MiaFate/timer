//go:build darwin

package desktop

import (
	"timer/internal/domain"
)

type DarwinService struct{}

func NewWindowService() domain.WindowService {
	return &DarwinService{}
}

func (s *DarwinService) GetOpenWindows() ([]domain.AppInfo, error) {
	return []domain.AppInfo{}, nil
}

func (s *DarwinService) GetForegroundWindowPID() (uint32, error) {
	return 0, nil
}
