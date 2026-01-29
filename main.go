package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"

	// Layers
	delivery "timer/internal/delivery/http"
	"timer/internal/infrastructure/desktop"
	"timer/internal/usecase"

	"github.com/gin-gonic/gin"
	webview "github.com/webview/webview_go"
)

//go:embed frontend/dist/*
var frontendFS embed.FS

func main() {
	// 1. Infrastructure (Cross-Platform)
	winService := desktop.NewWindowService()

	// 2. UseCase
	trackerUC := usecase.NewTrackerUseCase(winService)

	// 3. Delivery
	h := delivery.NewHandler(trackerUC)

	// Set Gin to release mode
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	r.SetTrustedProxies(nil)

	// CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	r.GET("/api/apps", h.GetApps)
	r.POST("/api/track", h.StartTracking)
	r.POST("/api/stop", h.StopTracking)
	r.GET("/api/status", h.GetStatus)

	// Serve Frontend
	distFS, err := fs.Sub(frontendFS, "frontend/dist")
	if err != nil {
		log.Fatal("Error serving frontend:", err)
	}
	r.NoRoute(gin.WrapH(http.FileServer(http.FS(distFS))))

	log.Println("Server starting on :8080")

	// Start server in goroutine
	go r.Run(":8080")

	// Launch Native Window
	w := webview.New(true)
	defer w.Destroy()
	w.SetTitle("Focus Tracker")
	w.SetSize(900, 700, webview.HintNone)
	w.Navigate("http://localhost:8080")

	// Force Icon using infrastructure helper
	desktop.SetWindowIcon(w.Window())

	w.Run()
}
