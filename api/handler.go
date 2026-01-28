package api

import (
	"net/http"
	"timer/platform"
	"timer/tracker"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Tracker *tracker.Tracker
}

func NewHandler(t *tracker.Tracker) *Handler {
	return &Handler{Tracker: t}
}

func (h *Handler) GetApps(c *gin.Context) {
	apps, err := platform.GetOpenWindows()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, apps)
}

func (h *Handler) StartTracking(c *gin.Context) {
	var req struct {
		PID uint32 `json:"pid"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	h.Tracker.Start(req.PID)
	c.JSON(http.StatusOK, gin.H{"status": "tracking started", "pid": req.PID})
}

func (h *Handler) StopTracking(c *gin.Context) {
	h.Tracker.Stop()
	c.JSON(http.StatusOK, gin.H{"status": "tracking stopped"})
}

func (h *Handler) GetStatus(c *gin.Context) {
	tracking, duration, pid, isActive := h.Tracker.GetStatus()
	c.JSON(http.StatusOK, gin.H{
		"tracking":  tracking,
		"seconds":   duration.Seconds(),
		"pid":       pid,
		"is_active": isActive,
	})
}
