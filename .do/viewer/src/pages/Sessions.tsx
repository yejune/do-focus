import { useEffect, useState } from 'react'
import { api, type Session } from '../api/client'

export default function Sessions() {
  const [sessions, setSessions] = useState<Session[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [selectedSession, setSelectedSession] = useState<Session | null>(null)

  useEffect(() => {
    async function loadSessions() {
      setLoading(true)
      try {
        const data = await api.getSessions()
        setSessions(data || [])
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to load sessions')
      } finally {
        setLoading(false)
      }
    }

    loadSessions()
  }, [])

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold text-gray-900">Sessions</h1>
        <p className="text-sm text-gray-500">{sessions.length} total</p>
      </div>

      {error && (
        <div className="bg-red-50 border border-red-200 rounded-lg p-4 text-red-700">
          {error}
        </div>
      )}

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Sessions List */}
        <div className="lg:col-span-2 bg-white rounded-lg shadow">
          <div className="p-4 border-b">
            <h2 className="font-semibold">All Sessions</h2>
          </div>
          <div className="divide-y max-h-[600px] overflow-y-auto">
            {loading ? (
              [...Array(5)].map((_, i) => (
                <div key={i} className="p-4 animate-pulse">
                  <div className="h-4 bg-gray-200 rounded w-1/3 mb-2" />
                  <div className="h-3 bg-gray-100 rounded w-1/2" />
                </div>
              ))
            ) : sessions.length === 0 ? (
              <div className="p-8 text-center text-gray-500">
                No sessions found
              </div>
            ) : (
              sessions.map((session) => (
                <button
                  key={session.id}
                  onClick={() => setSelectedSession(session)}
                  className={`w-full p-4 text-left hover:bg-gray-50 transition-colors ${
                    selectedSession?.id === session.id ? 'bg-primary-50' : ''
                  }`}
                >
                  <div className="flex items-center justify-between">
                    <div>
                      <p className="font-medium text-sm">
                        {session.project_path?.split('/').slice(-2).join('/') || session.id || 'Unknown'}
                      </p>
                      <p className="text-xs text-gray-500 mt-1">
                        {new Date(session.started_at).toLocaleString('ko-KR')}
                      </p>
                    </div>
                    <span className={`text-xs px-2 py-1 rounded ${
                      session.ended_at
                        ? 'bg-gray-100 text-gray-600'
                        : 'bg-green-100 text-green-700'
                    }`}>
                      {session.ended_at ? 'Ended' : 'Active'}
                    </span>
                  </div>
                  {session.summary && (
                    <p className="text-xs text-gray-600 mt-2 line-clamp-2">
                      {session.summary}
                    </p>
                  )}
                </button>
              ))
            )}
          </div>
        </div>

        {/* Session Details */}
        <div className="bg-white rounded-lg shadow">
          <div className="p-4 border-b">
            <h2 className="font-semibold">Session Details</h2>
          </div>
          <div className="p-4">
            {selectedSession ? (
              <div className="space-y-4">
                <div>
                  <label className="text-xs text-gray-500 uppercase tracking-wide">
                    Session ID
                  </label>
                  <p className="font-mono text-sm break-all">{selectedSession.id}</p>
                </div>
                <div>
                  <label className="text-xs text-gray-500 uppercase tracking-wide">
                    Project Path
                  </label>
                  <p className="text-sm break-all">{selectedSession.project_path}</p>
                </div>
                <div>
                  <label className="text-xs text-gray-500 uppercase tracking-wide">
                    Started At
                  </label>
                  <p className="text-sm">
                    {new Date(selectedSession.started_at).toLocaleString('ko-KR')}
                  </p>
                </div>
                {selectedSession.ended_at && (
                  <div>
                    <label className="text-xs text-gray-500 uppercase tracking-wide">
                      Ended At
                    </label>
                    <p className="text-sm">
                      {new Date(selectedSession.ended_at).toLocaleString('ko-KR')}
                    </p>
                  </div>
                )}
                {selectedSession.summary && (
                  <div>
                    <label className="text-xs text-gray-500 uppercase tracking-wide">
                      Summary
                    </label>
                    <p className="text-sm text-gray-700 whitespace-pre-wrap">
                      {selectedSession.summary}
                    </p>
                  </div>
                )}
              </div>
            ) : (
              <p className="text-gray-500 text-center py-8">
                Select a session to view details
              </p>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}
