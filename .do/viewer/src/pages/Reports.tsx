import { useEffect, useState } from 'react'

interface SessionSummary {
  id: number
  session_id: string
  type: string
  content: string
  created_at: string
}

export default function Reports() {
  const [summaries, setSummaries] = useState<SessionSummary[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [days, setDays] = useState(7)
  const [selectedSummary, setSelectedSummary] = useState<SessionSummary | null>(null)

  useEffect(() => {
    async function loadSummaries() {
      setLoading(true)
      try {
        const response = await fetch(`http://127.0.0.1:3778/api/summaries?days=${days}`)
        if (!response.ok) throw new Error('Failed to fetch')
        const data = await response.json()
        setSummaries(Array.isArray(data) ? data : [])
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to load summaries')
      } finally {
        setLoading(false)
      }
    }

    loadSummaries()
  }, [days])

  // Group summaries by date
  const groupedByDate = summaries.reduce((acc, s) => {
    const date = new Date(s.created_at).toLocaleDateString('ko-KR')
    if (!acc[date]) acc[date] = []
    acc[date].push(s)
    return acc
  }, {} as Record<string, SessionSummary[]>)

  const dates = Object.keys(groupedByDate).sort((a, b) =>
    new Date(b).getTime() - new Date(a).getTime()
  )

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold text-gray-900">Reports</h1>
        <div className="flex items-center gap-2">
          <label className="text-sm text-gray-500">Period:</label>
          <select
            value={days}
            onChange={(e) => setDays(Number(e.target.value))}
            className="border rounded-lg px-3 py-1.5 text-sm focus:ring-2 focus:ring-primary-500"
          >
            <option value={7}>Last 7 days</option>
            <option value={14}>Last 14 days</option>
            <option value={30}>Last 30 days</option>
          </select>
        </div>
      </div>

      {error && (
        <div className="bg-red-50 border border-red-200 rounded-lg p-4 text-red-700">
          {error}
        </div>
      )}

      {/* Summary Stats */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div className="bg-white rounded-lg shadow p-6">
          <p className="text-sm text-gray-500">Total Summaries</p>
          {loading ? (
            <div className="animate-pulse h-8 w-16 bg-gray-200 rounded mt-1" />
          ) : (
            <p className="text-2xl font-bold text-gray-900">{summaries.length}</p>
          )}
        </div>
        <div className="bg-white rounded-lg shadow p-6">
          <p className="text-sm text-gray-500">Active Days</p>
          {loading ? (
            <div className="animate-pulse h-8 w-16 bg-gray-200 rounded mt-1" />
          ) : (
            <p className="text-2xl font-bold text-gray-900">{dates.length}</p>
          )}
        </div>
        <div className="bg-white rounded-lg shadow p-6">
          <p className="text-sm text-gray-500">Unique Sessions</p>
          {loading ? (
            <div className="animate-pulse h-8 w-16 bg-gray-200 rounded mt-1" />
          ) : (
            <p className="text-2xl font-bold text-gray-900">
              {new Set(summaries.map(s => s.session_id)).size}
            </p>
          )}
        </div>
      </div>

      {/* Session Summaries */}
      <div className="bg-white rounded-lg shadow">
        <div className="p-4 border-b">
          <h2 className="font-semibold">Session Summaries</h2>
        </div>
        <div className="divide-y max-h-[600px] overflow-y-auto">
          {loading ? (
            [...Array(3)].map((_, i) => (
              <div key={i} className="p-4 animate-pulse">
                <div className="h-4 bg-gray-200 rounded w-1/4 mb-2" />
                <div className="h-3 bg-gray-100 rounded w-3/4" />
              </div>
            ))
          ) : summaries.length === 0 ? (
            <div className="p-8 text-center text-gray-500">
              No summaries for selected period
            </div>
          ) : (
            dates.map((date) => (
              <div key={date} className="p-4">
                <div className="flex items-center justify-between mb-3">
                  <p className="font-medium text-gray-900">{date}</p>
                  <span className="text-xs text-gray-500">
                    {groupedByDate[date].length} summaries
                  </span>
                </div>
                <div className="space-y-3">
                  {groupedByDate[date].map((summary) => (
                    <button
                      key={summary.id}
                      onClick={() => setSelectedSummary(summary)}
                      className="w-full text-left bg-gray-50 rounded-lg p-3 hover:bg-gray-100 transition-colors"
                    >
                      <div className="flex items-center justify-between mb-2">
                        <span className="text-xs font-mono text-gray-400">
                          {summary.session_id.slice(0, 8)}...
                        </span>
                        <span className="text-xs text-gray-400">
                          {new Date(summary.created_at).toLocaleTimeString('ko-KR')}
                        </span>
                      </div>
                      <div className="text-sm text-gray-700 whitespace-pre-wrap line-clamp-3">
                        {summary.content}
                      </div>
                      <span className="text-xs text-primary-600 mt-2 inline-block">
                        클릭하여 상세보기 →
                      </span>
                    </button>
                  ))}
                </div>
              </div>
            ))
          )}
        </div>
      </div>

      {/* Detail Modal */}
      {selectedSummary && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-lg shadow-xl max-w-3xl w-full max-h-[80vh] flex flex-col">
            <div className="p-4 border-b flex items-center justify-between">
              <div>
                <h3 className="font-semibold">Session Summary</h3>
                <p className="text-xs text-gray-500 font-mono">
                  {selectedSummary.session_id}
                </p>
              </div>
              <button
                onClick={() => setSelectedSummary(null)}
                className="text-gray-400 hover:text-gray-600"
              >
                <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>
            <div className="p-4 overflow-y-auto flex-1">
              <p className="text-xs text-gray-400 mb-2">
                {new Date(selectedSummary.created_at).toLocaleString('ko-KR')}
              </p>
              <div className="prose prose-sm max-w-none">
                <pre className="whitespace-pre-wrap text-sm text-gray-700 bg-gray-50 p-4 rounded-lg">
                  {selectedSummary.content}
                </pre>
              </div>
            </div>
            <div className="p-4 border-t flex justify-end">
              <button
                onClick={() => setSelectedSummary(null)}
                className="px-4 py-2 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200"
              >
                닫기
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
