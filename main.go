package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"syscall"
	"timer/api"
	"timer/tracker"
	"unsafe"

	"github.com/gin-gonic/gin"
	webview "github.com/webview/webview_go"
)

//go:embed frontend/dist/*
var frontendFS embed.FS

var (
	user32                = syscall.NewLazyDLL("user32.dll")
	kernel32              = syscall.NewLazyDLL("kernel32.dll")
	shell32               = syscall.NewLazyDLL("shell32.dll")
	procSendMessage       = user32.NewProc("SendMessageW")
	procLoadImage         = user32.NewProc("LoadImageW")
	procGetModuleHandle   = kernel32.NewProc("GetModuleHandleW")
	procGetModuleFileName = kernel32.NewProc("GetModuleFileNameW")
	procExtractIcon       = shell32.NewProc("ExtractIconW")
)

const (
	WM_SETICON     = 0x0080
	ICON_SMALL     = 0
	ICON_BIG       = 1
	IMAGE_ICON     = 1
	LR_DEFAULTSIZE = 0x0040
)

func main() {
	t := tracker.NewTracker()
	h := api.NewHandler(t)

	// Set Gin to release mode for cleaner output
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	// Fix: Disable trusted proxies warning
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
		log.Fatal(err)
	}
	r.NoRoute(gin.WrapH(http.FileServer(http.FS(distFS))))

	log.Println("Server starting on :8080")

	// Start server in goroutine
	go r.Run(":8080")

	// Launch Native Window
	w := webview.New(true) // true = debug mode
	defer w.Destroy()
	w.SetTitle("Focus Tracker")
	w.SetSize(900, 700, webview.HintNone)
	w.Navigate("http://localhost:8080")

	// Force Icon Load (Win32)
	hwnd := w.Window()
	setWindowIcon(hwnd)

	w.Run()
}

func setWindowIcon(hwnd unsafe.Pointer) {
	// robust strategy: Extract the icon from the executable file itself.
	// This grabs specific index 0 (the main icon).

	// 1. Get Path to Exe
	buf := make([]uint16, 260)
	procGetModuleFileName.Call(0, uintptr(unsafe.Pointer(&buf[0])), 260)

	// 2. Extract Icon
	iconHandle, _, _ := procExtractIcon.Call(
		0,
		uintptr(unsafe.Pointer(&buf[0])),
		0, // Index 0
	)

	if iconHandle == 0 || iconHandle == 1 {
		return
	}

	// 3. Set Icon
	procSendMessage.Call(uintptr(hwnd), WM_SETICON, ICON_BIG, iconHandle)
	procSendMessage.Call(uintptr(hwnd), WM_SETICON, ICON_SMALL, iconHandle)
}
