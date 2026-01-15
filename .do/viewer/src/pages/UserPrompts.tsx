import { useEffect, useState } from 'react'

interface UserPrompt {
  id: number
  session_id: string
  prompt_number: number
  prompt_text: string
  created_at: string
}

export default function UserPrompts() {
  const [prompts, setPrompts] = useState<UserPrompt[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    async function loadPrompts() {
      setLoading(true)
      try {
        // Note: This endpoint may need to be added to the worker
        const response = await fetch('http://127.0.0.1:3778/api/prompts?limit=100')
        if (!response.ok) throw new Error(`HTTP ${response.status}`)
        const data = await response.json()
        setPrompts(Array.isArray(data) ? data : [])
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to load prompts')
      } finally {
        setLoading(false)
      }
    }

    loadPrompts()
  }, [])

  // Group by session
  const groupedBySession = prompts.reduce((acc, p) => {
    if (!acc[p.session_id]) acc[p.session_id] = []
    acc[p.session_id].push(p)
    return acc
  }, {} as Record<string, UserPrompt[]>)

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold text-gray-900">User Prompts</h1>
        <p className="text-sm text-gray-500">{prompts.length} total</p>
      </div>

      {error && (
        <div className="bg-red-50 border border-red-200 rounded-lg p-4 text-red-700">
          {error}
        </div>
      )}

      {/* Stats */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div className="bg-white rounded-lg shadow p-6">
          <p className="text-sm text-gray-500">Total Prompts</p>
          {loading ? (
            <div className="animate-pulse h-8 w-16 bg-gray-200 rounded mt-1" />
          ) : (
            <p className="text-2xl font-bold text-gray-900">{prompts.length}</p>
          )}
        </div>
        <div className="bg-white rounded-lg shadow p-6">
          <p className="text-sm text-gray-500">Sessions with Prompts</p>
          {loading ? (
            <div className="animate-pulse h-8 w-16 bg-gray-200 rounded mt-1" />
          ) : (
            <p className="text-2xl font-bold text-gray-900">
              {Object.keys(groupedBySession).length}
            </p>
          )}
        </div>
      </div>

      {/* Prompts List */}
      <div className="bg-white rounded-lg shadow">
        <div className="p-4 border-b">
          <h2 className="font-semibold">Recent Prompts</h2>
        </div>
        <div className="divide-y max-h-[600px] overflow-y-auto">
          {loading ? (
            [...Array(5)].map((_, i) => (
              <div key={i} className="p-4 animate-pulse">
                <div className="h-4 bg-gray-200 rounded w-1/3 mb-2" />
                <div className="h-3 bg-gray-100 rounded w-2/3" />
              </div>
            ))
          ) : prompts.length === 0 ? (
            <div className="p-8 text-center text-gray-500">
              No prompts recorded yet
            </div>
          ) : (
            prompts.map((prompt) => (
              <div key={prompt.id} className="p-4 hover:bg-gray-50">
                <div className="flex items-center justify-between mb-2">
                  <span className="text-xs font-mono text-gray-400">
                    {prompt.session_id.slice(0, 8)}... #{prompt.prompt_number || prompt.id}
                  </span>
                  <span className="text-xs text-gray-400">
                    {new Date(prompt.created_at).toLocaleString('ko-KR')}
                  </span>
                </div>
                <p className="text-sm text-gray-700 line-clamp-3">
                  {prompt.prompt_text}
                </p>
              </div>
            ))
          )}
        </div>
      </div>
    </div>
  )
}
