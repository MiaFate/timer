package tracker

import (
	"fmt"
	"sync"
	"time"
	"timer/platform"
)

type Tracker struct {
	TargetPID       uint32
	AccumulatedTime time.Duration
	IsTracking      bool
	IsTargetActive  bool
	mu              sync.Mutex
	stopCh          chan struct{}
}

func NewTracker() *Tracker {
	return &Tracker{
		stopCh: make(chan struct{}),
	}
}

func (t *Tracker) Start(pid uint32) {
	t.mu.Lock()
	if t.IsTracking {
		t.Stop() // Stop previous if any
	}
	t.TargetPID = pid
	t.IsTracking = true
	t.IsTargetActive = false
	t.stopCh = make(chan struct{})
	t.mu.Unlock()

	go t.loop()
}

func (t *Tracker) Stop() {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.IsTracking {
		close(t.stopCh)
		t.IsTracking = false
		t.IsTargetActive = false
		t.AccumulatedTime = 0
	}
}

func (t *Tracker) loop() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-t.stopCh:
			return
		case <-ticker.C:
			activePID, err := platform.GetForegroundWindowPID()
			if err != nil {
				fmt.Println("Error getting foreground window:", err)
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

func (t *Tracker) GetStatus() (bool, time.Duration, uint32, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.IsTracking, t.AccumulatedTime, t.TargetPID, t.IsTargetActive
}
