package usecase

import (
	"fmt"
	"sync"
	"time"
	"timer/internal/domain"
)

type TrackerUseCase struct {
	TargetPID       uint32
	AccumulatedTime time.Duration
	IsTracking      bool
	IsTargetActive  bool
	mu              sync.Mutex
	stopCh          chan struct{}

	// Dependency
	winService domain.WindowService
}

func NewTrackerUseCase(ws domain.WindowService) *TrackerUseCase {
	return &TrackerUseCase{
		stopCh:     make(chan struct{}),
		winService: ws,
	}
}

func (t *TrackerUseCase) Start(pid uint32) {
	t.mu.Lock()
	if t.IsTracking {
		t.Stop()
	}
	t.TargetPID = pid
	t.IsTracking = true
	t.IsTargetActive = false
	t.stopCh = make(chan struct{})
	t.mu.Unlock()

	go t.loop()
}

func (t *TrackerUseCase) Stop() {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.IsTracking {
		close(t.stopCh)
		t.IsTracking = false
		t.IsTargetActive = false
		t.AccumulatedTime = 0
	}
}

func (t *TrackerUseCase) loop() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-t.stopCh:
			return
		case <-ticker.C:
			activePID, err := t.winService.GetForegroundWindowPID()
			if err != nil {
				fmt.Println("Error checking active window:", err)
				continue
			}

			t.mu.Lock()
			if activePID == t.TargetPID {
				t.AccumulatedTime += 1 * time.Second
				t.IsTargetActive = true
			} else {
				t.IsTargetActive = false
			}
			t.mu.Unlock()
		}
	}
}

func (t *TrackerUseCase) GetStatus() (bool, time.Duration, uint32, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.IsTracking, t.AccumulatedTime, t.TargetPID, t.IsTargetActive
}

func (t *TrackerUseCase) GetOpenApps() ([]domain.AppInfo, error) {
	// This could be in a separate AppUseCase, but fitting here for simplicity for now
	return t.winService.GetOpenWindows()
}
