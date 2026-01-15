import { useEffect, useState } from 'react'
import { api, type Summary } from '../api/client'

export default function Reports() {
  const [summaries, setSummaries] = useState<Summary[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [days, setDays] = useState(7)

  useEffect(() => {
    async function loadSummaries() {
      setLoading(true)
      try {
        const data = await api.getSummaries(days)
        setSummaries(data || [])
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to load summaries')
      } finally {
        setLoading(false)
      }
    }

    loadSummaries()
  }, [days])

  const totalSessions = summaries.reduce((sum, s) => sum + s.sessions_count, 0)
  const totalObservations = summaries.reduce((sum, s) => sum + s.observations_count, 0)

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
          <p className="text-sm text-gray-500">Total Sessions</p>
          {loading ? (
            <div className="animate-pulse h-8 w-16 bg-gray-200 rounded mt-1" />
          ) : (
            <p className="text-2xl font-bold text-gray-900">{totalSessions}</p>
          )}
        </div>
        <div className="bg-white rounded-lg shadow p-6">
          <p className="text-sm text-gray-500">Total Observations</p>
          {loading ? (
            <div className="animate-pulse h-8 w-16 bg-gray-200 rounded mt-1" />
          ) : (
            <p className="text-2xl font-bold text-gray-900">{totalObservations}</p>
          )}
        </div>
        <div className="bg-white rounded-lg shadow p-6">
          <p className="text-sm text-gray-500">Active Days</p>
          {loading ? (
            <div className="animate-pulse h-8 w-16 bg-gray-200 rounded mt-1" />
          ) : (
            <p className="text-2xl font-bold text-gray-900">
              {summaries.filter(s => s.sessions_count > 0).length}
            </p>
          )}
        </div>
      </div>

      {/* Daily Summaries */}
      <div className="bg-white rounded-lg shadow">
        <div className="p-4 border-b">
          <h2 className="font-semibold">Daily Summaries</h2>
        </div>
        <div className="divide-y">
          {loading ? (
            [...Array(5)].map((_, i) => (
              <div key={i} className="p-4 animate-pulse">
                <div className="h-4 bg-gray-200 rounded w-1/4 mb-2" />
                <div className="h-3 bg-gray-100 rounded w-1/2" />
              </div>
            ))
          ) : summaries.length === 0 ? (
            <div className="p-8 text-center text-gray-500">
              No data for selected period
            </div>
          ) : (
            summaries.map((summary) => (
              <div key={summary.date} className="p-4">
                <div className="flex items-center justify-between mb-2">
                  <p className="font-medium">
                    {new Date(summary.date).toLocaleDateString('ko-KR', {
                      weekday: 'short',
                      month: 'short',
                      day: 'numeric',
                    })}
                  </p>
                  <div className="flex items-center gap-4 text-sm text-gray-500">
                    <span>{summary.sessions_count} sessions</span>
                    <span>{summary.observations_count} observations</span>
                  </div>
                </div>
                {summary.highlights.length > 0 && (
                  <ul className="space-y-1 mt-2">
                    {summary.highlights.map((highlight, i) => (
                      <li key={i} className="text-sm text-gray-600 flex items-start gap-2">
                        <span className="text-primary-500 mt-1">-</span>
                        {highlight}
                      </li>
                    ))}
                  </ul>
                )}
                {summary.sessions_count === 0 && (
                  <p className="text-sm text-gray-400 italic">No activity</p>
                )}
              </div>
            ))
          )}
        </div>
      </div>
    </div>
  )
}
