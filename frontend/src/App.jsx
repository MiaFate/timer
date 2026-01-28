import { useState, useEffect } from 'react'

function App() {
  const [apps, setApps] = useState([])
  const [status, setStatus] = useState({ tracking: false, seconds: 0, pid: 0, is_active: false })
  const [loading, setLoading] = useState(false)
  const [searchTerm, setSearchTerm] = useState("")

  useEffect(() => {
    fetchApps()
    const interval = setInterval(fetchStatus, 1000)
    return () => clearInterval(interval)
  }, [])

  const fetchApps = async () => {
    try {
      const res = await fetch('http://localhost:8080/api/apps')
      const data = await res.json()
      setApps(data || [])
    } catch (e) {
      console.error(e)
    }
  }

  const fetchStatus = async () => {
    try {
      const res = await fetch('http://localhost:8080/api/status')
      const data = await res.json()
      setStatus(data)
    } catch (e) {
      console.error(e)
    }
  }

  const startTracking = async (pid) => {
    setLoading(true)
    try {
      await fetch('http://localhost:8080/api/track', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ pid }),
      })
      await fetchStatus()
    } finally {
      setLoading(false)
    }
  }

  const stopTracking = async () => {
    setLoading(true)
    try {
      await fetch('http://localhost:8080/api/stop', { method: 'POST' })
      await fetchStatus()
    } finally {
      setLoading(false)
    }
  }

  const formatTime = (seconds) => {
    const h = Math.floor(seconds / 3600)
    const m = Math.floor((seconds % 3600) / 60)
    const s = Math.floor(seconds % 60)
    return `${h}h ${m}m ${s}s`
  }

  const filteredApps = apps.filter(app =>
    app.title.toLowerCase().includes(searchTerm.toLowerCase())
  )

  return (
    <div className="min-h-screen bg-slate-900 text-slate-100 p-8 flex flex-col items-center font-sans">
      <h1 className="text-4xl font-bold mb-8 text-blue-400 drop-shadow-lg">
        Focus Tracker
      </h1>

      {status.tracking ? (
        <div className="flex flex-col items-center w-full max-w-md animate-in fade-in zoom-in duration-300 gap-6">
          <div className="bg-[#43484b] p-10 rounded-2xl shadow-2xl flex flex-col items-center w-full border border-slate-700">
            <div className="text-xl text-slate-400 mb-4">
              {status.is_active ? "You are tracking" : "Paused"}
            </div>
            <div className="text-6xl font-mono font-bold text-[#20bccd] tabular-nums">
              {formatTime(status.seconds)}
            </div>
          </div>

          <button
            onClick={stopTracking}
            disabled={loading}
            className="px-8 py-3 bg-red-500 hover:bg-red-600 rounded-lg font-semibold transition-all shadow-lg shadow-red-500/20 active:scale-95 cursor-pointer disabled:opacity-50 w-full"
          >
            Stop Tracking
          </button>
        </div>
      ) : (
        <div className="w-full max-w-2xl">
          <div className="flex justify-between items-center mb-6">
            <h2 className="text-2xl font-semibold text-slate-200">Select Application</h2>
            <button
              onClick={fetchApps}
              className="text-sm text-slate-400 hover:text-white transition-colors cursor-pointer"
            >
              Refresh List
            </button>
          </div>

          <input
            type="text"
            placeholder="Search applications..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            className="w-full p-3 mb-6 bg-slate-800 border border-slate-700 rounded-lg text-slate-200 focus:outline-none focus:border-blue-500 transition-colors"
          />

          <div className="grid gap-3">
            {filteredApps.length === 0 && (
              <div className="text-center text-slate-500 py-10">No applications found.</div>
            )}
            {filteredApps.map((app) => (
              <div
                key={app.pid}
                onClick={() => startTracking(app.pid)}
                className="bg-slate-800/50 hover:bg-slate-800 p-4 rounded-xl border border-slate-700/50 hover:border-blue-500/50 transition-all cursor-pointer flex justify-between items-center group"
              >
                <span className="font-medium truncate pr-4 text-slate-300 group-hover:text-white">
                  {app.title}
                </span>
                <span className="text-xs text-slate-600 font-mono bg-slate-900 px-2 py-1 rounded">
                  PID: {app.pid}
                </span>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  )
}

export default App
