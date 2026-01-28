# Focus Tracker App
## Features
- **Native Window**: Clean desktop interface using WebView.
- **Smart Filtering**: Deduplicated application list.
- **Search**: Real-time filtering.
- **Visuals**: Dark mode, custom timer colors, paused state.
- **Icon**: Official Go Gopher logo on executable and taskbar.
## Usage
1.  Run `gogogott.exe`.
2.  Type in the search box to find an app.
3.  Click to track. The timer counts when that window is active.
4.  Stop tracking to reset the timer.
## Technical Details
- **Stack**: Go (Backend) + React/Vite (Frontend).
- **Architecture**: Single executable with embedded assets.
- **Icon**: Embedded using `rsrc`.
## How to Build (From Source)
```powershell
# 1. Build Frontend
cd frontend
npm install
npm run build
cd ..
# 2. Embed Icon
go install github.com/akavel/rsrc@latest
rsrc -ico icon.ico
# 3. Build Backend (Hidden Console)
go get github.com/webview/webview_go
go build -ldflags="-H windowsgui" -o gogogott.exe .