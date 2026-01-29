package http

import (
	"net/http"
	"timer/internal/usecase"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	TrackerUseCase *usecase.TrackerUseCase
}

func NewHandler(uc *usecase.TrackerUseCase) *Handler {
	return &Handler{TrackerUseCase: uc}
}

func (h *Handler) GetApps(c *gin.Context) {
	apps, err := h.TrackerUseCase.GetOpenApps()
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

	h.TrackerUseCase.Start(req.PID)
	c.JSON(http.StatusOK, gin.H{"status": "tracking started", "pid": req.PID})
}

func (h *Handler) StopTracking(c *gin.Context) {
	h.TrackerUseCase.Stop()
	c.JSON(http.StatusOK, gin.H{"status": "tracking stopped"})
}

func (h *Handler) GetStatus(c *gin.Context) {
	tracking, duration, pid, isActive := h.TrackerUseCase.GetStatus()
	c.JSON(http.StatusOK, gin.H{
		"tracking":  tracking,
		"seconds":   duration.Seconds(),
		"pid":       pid,
		"is_active": isActive,
	})
}
